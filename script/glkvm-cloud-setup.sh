#!/bin/bash
echo "GLKVM cloud is building..."
sed -i "s/#\$nrconf{restart} = 'i';/\$nrconf{restart} = 'a';/" /etc/needrestart/needrestart.conf

GLKVM_DIR="$PWD/glkvm_cloud"

if [ ! -d "$GLKVM_DIR" ]; then
    mkdir -p "$GLKVM_DIR"
    echo "Created directory: $GLKVM_DIR"
else
    echo "Directory already exists: $GLKVM_DIR"
fi

IMAGE_URL="https://aw-test.gl-inet.cn/server-node/selfhost/glkvmcloud.tar"
IMAGE_PATH="$GLKVM_DIR/glkvm-cloud.tar"

echo "Downloading Docker image from: $IMAGE_URL"
curl -L -o "$IMAGE_PATH" "$IMAGE_URL"

if [ $? -ne 0 ]; then
    echo "❌ Failed to download image. Please check network or URL."
    exit 1
fi

echo "Downloaded image to: $IMAGE_PATH"

echo "Importing Docker image..."
docker load -i "$IMAGE_PATH"

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
version: '3.8'

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

echo "GLKVM Cloud has been successfully initialized: $GLKVM_DIR"
echo "🔐 GLKVM Cloud Web access password: $PASSWORD"
echo "Next steps: configure your domain and SSL certificate (refer to the domain setup guide), then start your GLKVM Cloud with 'docker-compose up -d'."
