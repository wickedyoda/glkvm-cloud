<template>
  <el-dialog v-model="scriptModel" :title="$t('Please run script in device shell to connected cloud')" width="800">
    <div class="script-header">
      <el-button icon="DocumentCopy" type="danger" size="small" @click="copyScript">{{ $t('Copy Script') }}</el-button>
    </div>

    <pre class="script-box"><code ref="codeRef" class="language-bash">{{ scriptContent }}</code></pre>

    <template #footer>
      <el-button type="primary" @click="closeScriptDialog">{{ $t('Close') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import Prism from 'prismjs'
import 'prismjs/themes/prism-tomorrow.css'
import 'prismjs/components/prism-bash.js'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const scriptModel = ref(false)
const codeRef = ref(null)

const scriptContent = ref('')

const generateScript = (hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword) => {
  return `#!/bin/sh

HOSTNAME="${hostname}"
PORT="${port}"
TOKEN="${token}"
WEBRTC_IP="${webrtcIP}"
WEBRTC_PORT="${webrtcPort}"
WEBRTC_USERNAME="${webrtcUsername}"
WEBRTC_PASSWORD="${webrtcPassword}"

TARGET_DIR="/etc/kvmd/user/scripts"
SCRIPT_FILE="$TARGET_DIR/S01selfCloud"
WATCHDOG_SCRIPT="$TARGET_DIR/rtty-loop.sh"

# 1. Create directory
mkdir -p "$TARGET_DIR"

# 2. Write S01selfCloud script
cat <<'EOF' > "$SCRIPT_FILE"
#!/bin/sh

HOSTNAME="${hostname}"
PORT="${port}"
TOKEN="${token}"
WEBRTC_IP="${webrtcIP}"
WEBRTC_PORT="${webrtcPort}"
WEBRTC_USERNAME="${webrtcUsername}"
WEBRTC_PASSWORD="${webrtcPassword}"

TARGET_DIR="/etc/kvmd/user/scripts"
SCRIPT_FILE="$TARGET_DIR/S01selfCloud"
WATCHDOG_SCRIPT="$TARGET_DIR/rtty-loop.sh"

start() {
    if pgrep -f "$WATCHDOG_SCRIPT" > /dev/null; then
        echo "selfHosted cloud is already started"
        return 0
    fi

    echo "starting self hosted cloud"

    # Write TURN config
    cat <<TURNCONF > /tmp/turnserver.json
{
    "username": "$WEBRTC_USERNAME",
    "ttl": 864000,
    "password": "$WEBRTC_PASSWORD",
    "uris": [
        "turn:$WEBRTC_IP:$WEBRTC_PORT?transport=udp",
        "turn:$WEBRTC_IP:$WEBRTC_PORT?transport=tcp"
    ]
}
TURNCONF

    device_id=$(cat /proc/gl-hw-info/device_ddns)
    device_mac=$(cat /proc/gl-hw-info/device_mac)

    mkdir -p $(dirname "$WATCHDOG_SCRIPT")

    cat <<RTTY > "$WATCHDOG_SCRIPT"
#!/bin/sh
device_id=$(cat /proc/gl-hw-info/device_ddns)
device_mac=$(cat /proc/gl-hw-info/device_mac)
while true; do
    if ! pgrep -f "rtty.*-d $device_mac" > /dev/null; then
        echo "rtty not running, starting..."
        rtty -sx -T 2 -I "$device_id" -h $HOSTNAME$PORT -t $TOKEN -d "$device_mac" &
    fi
    sleep 5
done
RTTY

    chmod +x "$WATCHDOG_SCRIPT"
    nohup "$WATCHDOG_SCRIPT" >> /tmp/rtty.log 2>&1 &
    echo "started selfHosted cloud"
}

stop() {
    echo "stopping selfHosted cloud"
    device_mac=$(cat /proc/gl-hw-info/device_mac)
    pkill -f "$WATCHDOG_SCRIPT"
    pkill -f "rtty.*-d $device_mac"
    echo "stopped selfHosted cloud"
}

restart() {
    stop
    sleep 1
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart|reload)
        restart
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
esac

exit 0
EOF

# 3. Add execution permissions
chmod +x "$SCRIPT_FILE"

# 4. Execute restart once
"$SCRIPT_FILE" restart

`
}

const showScriptDialog = async () => {
  try {
    const res = await axios.get('/get/scriptInfo')
    const { hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword } = res.data
    scriptContent.value = generateScript(hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword)

    scriptModel.value = true
  } catch (err) {
    ElMessage.error(t('Get script info failed'))
  }
}

const closeScriptDialog = () => {
  scriptModel.value = false
}

const copyScript = async () => {
  try {
    await navigator.clipboard.writeText(scriptContent.value)
    ElMessage.success(t('Copied'))
  } catch (e) {
    ElMessage.error(t('Copied failed'))
  }
}

watch(scriptModel, (val) => {
  if (val) {
    nextTick(() => {
      Prism.highlightElement(codeRef.value)
    })
  }
})

defineExpose({
  showScriptDialog,
})
</script>

<style scoped>
.script-box {
  border-radius: 6px;
  background: #f5f5f5;
  color: #333;         
  padding: 16px;
  font-family: monospace;
  font-size: 14px;
  overflow-x: auto;
  overflow-y: auto;
  line-height: 1.5;
  max-height: 400px; /* 限制高度，防止撑高对话框 */
}

.script-box::-webkit-scrollbar {
  height: 4px;  /* 横向滚动条高度 */
  width: 4px;   /* 纵向滚动条宽度 */
}

.script-box::-webkit-scrollbar-thumb {
  background: #aaa;
  border-radius: 4px;
}

.script-box::-webkit-scrollbar-thumb:hover {
  background: #888;
}

.script-header {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 10px;
}
</style>
