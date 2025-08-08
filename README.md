# Self-Deployed Lightweight Cloud KVM Remote Management Platform

Self-Deployed Lightweight Cloud is a lightweight KVM remote cloud platform tailored for individuals and small businesses. This project is developed based on [rttys](https://github.com/zhaojh329/rttys), specifically designed for users managing fewer than 50 KVM devices.

## Core Features

- 🚀 **Rapid Deployment** - Quick self-deployment with simple operations
- 🔒 **Data Security** - Private deployment with full data control
- 🎯 **Lightweight Design** - Optimized for small businesses and individual users
- 🌐 **Unified Management** - Centralized management and remote control via public servers

## Main Functions

- ✅ **User Authentication** - Secure login authentication system
- 📋 **Device Management** - Online device list monitoring
- 📜 **Script Deployment** - Convenient script-based device addition
- 🖥️ **Remote SSH** - Secure SSH remote connections
- 🎮 **Remote Control** - Graphical remote desktop control
- ⚡ **Batch Operations** - Batch command execution capabilities

## Feature Demonstrations

###  Add KVM Devices to the Lightweight Cloud 

- Copy script

![Add KVM Devices to the Lightweight Cloud](img/adddevice.png)

- Run the script in the device terminal

![Add KVM Devices to the Lightweight Cloud](img/rundevicescript.png)

- Devices connected to the cloud

![Add KVM Devices to the Lightweight Cloud](img/devicelist.png)


### Remote SSH Connection

![Remote SSH Screenshot](img/ssh.png)

### Remote Desktop Control

![Remote Control Screenshot](img/web.png)

## Supported Operating Systems

The following mainstream operating systems have been tested and verified:

### Debian Family

- Ubuntu 18.04 / 20.04 / 22.04 / 24.04
- Debian 11 / 12

### Red Hat Family

- AlmaLinux 8 / 9
- Rocky Linux 8 / 9
- CentOS Stream 9

## System Requirements

### Server Configuration Requirements

| Component    | Minimum Requirement |
| ------------ | ------------------- |
| CPU          | 2 cores or above    |
| Memory       | ≥ 2 GB              |
| Storage      | ≥ 40 GB             |
| KVM Firmware | ≥ v1.5.0            |

### Recommended Scale

- **Recommended Device Count**: Under 50 devices
- **Use Cases**: Individual users, small businesses, data security-sensitive environments

## Self-Hosting Guide

### 🔐 Cloud Security Group Settings

If your server provider uses a **cloud security group** (e.g., AWS, Aliyun, etc.),
 please make sure the following ports are **open and not occupied by other services**:

| Port  | Protocol | Purpose                        |
| ----- | -------- | ------------------------------ |
| 443   | TCP      | Web UI access                  |
| 10443 | TCP      | WebSocket proxy                |
| 5912  | TCP      | Device connection              |
| 3478  | TCP/UDP  | TURN server for WebRTC support |

⚠️ **Important:**
 These ports will be **used by GLKVM Cloud**.
 Please ensure **no other applications or services** on your server are binding to these ports,
 otherwise the lightweight cloud platform may fail to start properly.

------
### 📦 Installation

Run the following command **as root** to install GLKVM Cloud:

```bash
bash <(curl -fsSL https://aw-test.gl-inet.cn/server-node/selfhost/glkvm-cloud-setup.sh)
```

### 🌐 Platform Access

Once the installation is complete, access the platform via:

```
https://<your_server_public_ip>
```

⚠️ **Note**: Accessing via IP address will trigger a **browser certificate warning**.
 To eliminate the warning, it's strongly recommended to configure a **custom domain** with a valid SSL certificate.

### 🔑 Web UI Login Password

The default login password for the Web UI will be displayed in the installation script output:

```
🔐 Please check the installation console for your web login password.
```

![](img\password.png)

⚠️ **Note**:
 If you **don’t mind the browser certificate warning**, and only need to **remotely access a single KVM device**,
 you can **skip** configuring a custom domain and SSL certificate.

You can still access the platform using the server’s **public IP** with HTTPS:

```
https://<your_server_public_ip>
```

However, for production use or to support **multiple device access via subdomains**,
 it is highly recommended to configure your own **wildcard SSL certificate** (see below).

------

### 📁 Default SSL Certificate Path

The default self-signed SSL certificates are located at:

```
~/glkvm_cloud/certificate
```

### 🌐 Recommended DNS Records

To enable full domain-based access, configure the following DNS records for your domain:

```
┌────────────┬──────┬────────────────────┬─────────────────────────────┐
│ Hostname   │ Type │     Value           │         Purpose             │
├────────────┼──────┼────────────────────┼─────────────────────────────┤
│ www        │  A   │ Your public IP      │ Web access to the platform │
│ *          │  A   │ Your public IP      │ Remote access to KVMs      │
└────────────┴──────┴────────────────────┴─────────────────────────────┘
```

### 🔧 Using a Custom SSL Certificate

To avoid browser warnings, replace the default certificates with your own **wildcard SSL certificate**
 that supports both:

- `*.your-domain.com` (for device access)
- `www.your-domain.com` (for platform access)

Replace the following files in:

```
~/glkvm_cloud/certificate
```

- `glkvm.cer`
- `glkvm.key`

⚠️ **Make sure the filenames remain unchanged.**

### 🔄 Restart Services After Certificate Replacement

After replacing the certificates, restart the GLKVM Cloud services to apply the changes:

```bash
cd ~/glkvm_cloud
docker-compose down
docker-compose up -d
```

Or, on systems with the Docker CLI plugin:

```bash
docker compose down
docker compose up -d
```

###  Domain-Based Access Example

Once everything is configured, you can access the platform via your domain:

```
https://www.your-domain.com
```