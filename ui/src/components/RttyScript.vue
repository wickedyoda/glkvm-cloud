<template>
  <el-dialog v-model="scriptModel" :title="'请在设备端执行脚本连接平台'" width="800">
    <div class="script-header">
      <el-button icon="DocumentCopy" type="danger" size="small" @click="copyScript">复制脚本</el-button>
    </div>

    <pre class="script-box"><code ref="codeRef" class="language-bash">{{ scriptContent }}</code></pre>

    <template #footer>
      <el-button type="primary" @click="closeScriptDialog">关闭</el-button>
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

const scriptModel = ref(false)
const codeRef = ref(null)

const scriptContent = ref('')

const generateScript = (hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword) => {
  return `#!/bin/sh

cat <<EOF > /tmp/turnserver.json
{
    "username": "${webrtcUsername}",
    "ttl": 864000,
    "password": "${webrtcPassword}",
    "uris": [
        "turn:${webrtcIP}:${webrtcPort}?transport=udp",
        "turn:${webrtcIP}:${webrtcPort}?transport=tcp"
    ]
}
EOF

# Configurable variables
LOGFILE="/tmp/rtty.log"
WATCHDOG_SCRIPT="/tmp/rtty-loop.sh"
LOCKFILE="/tmp/rtty_cloud.pid"

# Clean up old lock file
rm -f "$LOCKFILE"

# Read device information
device_id=$(cat /proc/gl-hw-info/device_ddns)
device_mac=$(cat /proc/gl-hw-info/device_mac)

# Kill existing rtty instances and previous watchdog process
pkill -f "rtty.*-d $device_mac"
pkill -f "$WATCHDOG_SCRIPT"

echo "Current KVM Device ID : $device_id"
echo "Current KVM Device MAC: $device_mac"

# Write the watchdog script content
cat <<EOF > "$WATCHDOG_SCRIPT"
#!/bin/sh

device_id=\$(cat /proc/gl-hw-info/device_ddns)
device_mac=\$(cat /proc/gl-hw-info/device_mac)

while true; do
  if ! pgrep -f "rtty.*-d \$device_mac" > /dev/null; then
    echo "[\$(date)] rtty not running, starting..."
    rtty -sx -T 2 -I "\$device_id" -h ${hostname}${port} -t ${token} -d "\$device_mac" &
  fi
  sleep 5
done
EOF

chmod +x "$WATCHDOG_SCRIPT"

# Run the watchdog script in the background
nohup "$WATCHDOG_SCRIPT" >> "$LOGFILE" 2>&1 &
echo $! > "$LOCKFILE"

echo "rtty watchdog started in background for device: $device_id"

`
}

const showScriptDialog = async () => {
  try {
    const res = await axios.get('/get/scriptInfo')
    const { hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword } = res.data
    scriptContent.value = generateScript(hostname, port, token, webrtcIP, webrtcPort, webrtcUsername, webrtcPassword)

    scriptModel.value = true
  } catch (err) {
    ElMessage.error('获取脚本信息失败')
  }
}

const closeScriptDialog = () => {
  scriptModel.value = false
}

const copyScript = async () => {
  try {
    await navigator.clipboard.writeText(scriptContent.value)
    ElMessage.success('已复制')
  } catch (e) {
    ElMessage.error('复制失败')
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
