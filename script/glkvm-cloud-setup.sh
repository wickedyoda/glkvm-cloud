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

echo ""
echo "✅ GLKVM Cloud has been successfully initialized at:"
echo "   $GLKVM_DIR"
echo ""
echo "📄 SSL certificates have been downloaded to:"
echo "   $GLKVM_DIR/certificate"
echo "   You may replace the default certificate files with your own:"
echo "     - Keep the filenames the same:"
echo "         glkvm.cer"
echo "         glkvm.key"
echo ""
echo "🚪 Please ensure the following ports are open (both in your firewall and cloud security group):"
echo "     - 443/TCP       (Web UI access)"
echo "     - 10443/TCP     (WebSocket proxy)"
echo "     - 5912/TCP      (Device connection)"
echo "     - 3478/TCP/UDP  (TURN server for WebRTC)"
echo ""
echo "🚀 Next steps:"
echo "   1. Point your domain to this server's IP address."
echo "   2. Replace the SSL certificate files if needed (see above)."
echo "   3. Start GLKVM Cloud services:"
echo "      cd $GLKVM_DIR && docker-compose up -d"
echo "         or"
echo "      cd $GLKVM_DIR && docker compose up -d"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔐 Web access password:"
echo ""
echo "     $PASSWORD"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"


