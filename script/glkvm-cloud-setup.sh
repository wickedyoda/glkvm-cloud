#!/bin/bash
echo "GLKVM cloud is building..."

# 初始化平台变量
PLATFORM="unknown"

# 识别系统类型
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS_ID=$ID
    OS_ID_LIKE=$ID_LIKE
else
    echo "Cannot determine OS. Exiting."
    exit 1
fi

# Debian/Ubuntu 系列
if [[ "$OS_ID" == "debian" || "$OS_ID" == "ubuntu" || "$OS_ID_LIKE" == *"debian"* ]]; then
    sed -i "s/#\$nrconf{restart} = 'i';/\$nrconf{restart} = 'a';/" /etc/needrestart/needrestart.conf
    PLATFORM="debian"
    echo "Detected Debian-based system: $PRETTY_NAME"

    apt update
    apt install -y docker.io docker-compose curl ufw

    # 防火墙规则
    ufw allow 443/tcp
    ufw allow 10443/tcp
    ufw allow 5912/tcp
    ufw allow 3478/tcp
    ufw allow 3478/udp
    echo "Firewall rules updated via UFW."

# RedHat/CentOS/AlmaLinux 系列
elif [[ "$OS_ID" == "centos" || "$OS_ID" == "rhel" || "$OS_ID" == "almalinux" || "$OS_ID" == "rocky" || "$OS_ID_LIKE" == *"rhel"* ]]; then
    PLATFORM="redhat"
    echo "Detected Red Hat-based system: $PRETTY_NAME"

    dnf makecache
    dnf install -y curl dnf-plugins-core

    # 添加 Docker 源
    dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

    # 安装 Docker
    dnf install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

    # 启用服务
    systemctl enable --now docker

    # 防火墙规则
    firewall-cmd --permanent --add-port=443/tcp
    firewall-cmd --permanent --add-port=10443/tcp
    firewall-cmd --permanent --add-port=5912/tcp
    firewall-cmd --permanent --add-port=3478/tcp
    firewall-cmd --permanent --add-port=3478/udp
    firewall-cmd --reload
    echo "Firewall rules updated via firewalld."

else
    echo "Unsupported OS: $PRETTY_NAME"
    exit 1
fi

# ✅ 输出平台标记（后续逻辑可用 $PLATFORM 判断）
echo "Platform detected: $PLATFORM"


GLKVM_DIR="$PWD/glkvm_cloud"

if [ ! -d "$GLKVM_DIR" ]; then
    mkdir -p "$GLKVM_DIR"
    echo "Created directory: $GLKVM_DIR"
else
    echo "Directory already exists: $GLKVM_DIR"
fi


BASE_DOMAIN="https://aw-test.gl-inet.cn"

IMAGE_URL="$BASE_DOMAIN/server-node/selfhost/glkvmcloud.tar"
IMAGE_PATH="$GLKVM_DIR/glkvm-cloud.tar"
echo "Downloading Docker image from: $IMAGE_URL"
curl -L -o "$IMAGE_PATH" "$IMAGE_URL"
if [ $? -ne 0 ]; then
    echo "❌ Failed to download image. Please check network or URL."
    exit 1
fi
echo "Downloaded image to: $IMAGE_PATH"


COTURN_URL="$BASE_DOMAIN/server-node/selfhost/glkvmcoturn.tar"
COTURN_PATH="$GLKVM_DIR/glkvm-coturn.tar"
echo "Downloading Docker image from: $COTURN_URL"
curl -L -o "$COTURN_PATH" "$COTURN_URL"
if [ $? -ne 0 ]; then
    echo "❌ Failed to download image. Please check network or URL."
    exit 1
fi
echo "Downloaded image to: $COTURN_PATH"

echo "Importing Docker image..."
docker load -i "$IMAGE_PATH"
docker load -i "$COTURN_PATH"

if [ $? -eq 0 ]; then
    echo "✅ Docker image imported successfully."
else
    echo "❌ Docker image import failed."
    exit 1
fi

cd "$GLKVM_DIR" || exit 1

get_public_ip() {
    ip=$(curl -s --max-time 5 https://api.ipify.org)
    if [[ -n "$ip" ]]; then
        echo "$ip"
        return 0
    fi

    ip=$(curl -s --max-time 5 https://ifconfig.me)
    if [[ -n "$ip" ]]; then
        echo "$ip"
        return 0
    fi

    return 1
}


PUBLIC_IP=$(get_public_ip)
if [[ -z "$PUBLIC_IP" ]]; then
    echo "❌ Failed to get public IP from both sources. Please check your network."
    exit 1
fi

echo "Detected public IP: $PUBLIC_IP"


generate_random_string() {
    tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32
}

TOKEN=$(generate_random_string)
PASSWORD=$(generate_random_string)
WEBRTC_USERNAME=$(generate_random_string)
WEBRTC_PASSWORD=$(generate_random_string)

mkdir -p certificate

CERT_DIR="$GLKVM_DIR/certificate"
CER_URL="$BASE_DOMAIN/server-node/selfhost/glkvm.cer"
KEY_URL="$BASE_DOMAIN/server-node/selfhost/glkvm.key"

curl -L -o "$CERT_DIR/glkvm.cer" "$CER_URL"
if [ $? -ne 0 ]; then
    echo "❌ Failed to download glkvm.cer"
    
fi

curl -L -o "$CERT_DIR/glkvm.key" "$KEY_URL"
if [ $? -ne 0 ]; then
    echo "❌ Failed to download glkvm.key"
  
fi

cat > rttys.conf <<EOF
# Authentication token for device connections
token: $TOKEN

# Web management password
password: $PASSWORD

# Webrtc
webrtc-ip: $PUBLIC_IP
webrtc-port: 3478
webrtc-username: $WEBRTC_USERNAME
webrtc-password: $WEBRTC_PASSWORD

addr-dev: :5912
addr-user: :443
addr-http-proxy: :10443
EOF

cat > docker-compose.yml <<EOF
version: "3.3"

services:
  rttys:
    container_name: glkvm_cloud
    image: glkvm:v1
    ports:
      - "443:443"
      - "10443:10443"
      - "5912:5912"
    volumes:
      - ./rttys.conf:/home/rttys.conf:ro
      - ./certificate/glkvm.cer:/home/certificate/glkvm_cer:ro
      - ./certificate/glkvm.key:/home/certificate/glkvm_key:ro
    command: ["-c", "/home/rttys.conf"]
    restart: always

  coturn:
    image: coturn/coturn:edge-alpine
    container_name: glkvm_coturn
    restart: always
    ports:
      - "3478:3478"
      - "3478:3478/udp"
      - "5349:5349"
      - "5349:5349/udp"
    volumes:
      - ./turnserver.conf:/etc/turnserver.conf:ro
    command: ["-c", "/etc/turnserver.conf"]
EOF

cat > turnserver.conf <<EOF
listening-port=3478
lt-cred-mech
user=$WEBRTC_USERNAME:$WEBRTC_PASSWORD
realm=glkvm
no-multicast-peers
allowed-peer-ip=0.0.0.0/0
EOF

if [ "$PLATFORM" = "debian" ]; then
    cd $GLKVM_DIR && docker-compose up -d
   
else
    cd $GLKVM_DIR && docker compose up -d
fi
echo ""
echo "✅ GLKVM Cloud has been successfully initialized at:"
echo "   $GLKVM_DIR"
echo ""
echo "   If your server provider enforces a cloud security group (e.g., on AWS, Aliyun, etc.),"
echo "   please ensure the following ports are allowed through:"
echo ""
echo "     - 443/TCP       (Web UI access)"
echo "     - 10443/TCP     (Device Web remote access)"
echo "     - 5912/TCP      (Device connection)"
echo "     - 3478/TCP/UDP  (TURN server for WebRTC)"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔐 Web Access:"
echo ""
echo "   🌐 You can now access the GLKVM Cloud platform via:"
echo "       https://$PUBLIC_IP"
echo ""
echo "   ⚠️  Note: Accessing via IP will trigger a browser certificate warning."
echo "       To remove this warning, please configure your own domain and SSL certificate."
echo ""
echo "   🔑 Web UI password:"
echo "       $PASSWORD"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

