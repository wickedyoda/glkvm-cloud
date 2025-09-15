# Quick Start

 This guide shows how to deploy **glkvm-cloud** using the provided Docker Compose environment template. 

1. **Clone the repository and prepare the environment template**

   ```bash
   git clone https://github.com/gl-inet/glkvm-cloud.git
   cd glkvm-cloud/docker-compose/
   cp .env.example .env
   ```

2. **Configure environment variables**

   Edit `.env` and update the required parameters: 

   - `RTTYS_TOKEN`: device connection token (leave empty to use the default)
   - `RTTYS_PASS`: web management password (leave empty to use the default **StrongP@ssw0rd**)
   - `TURN_USER` / `TURN_PASS`: coturn authentication credentials (leave empty to use the default)
   - `GLKVM_ACCESS_IP`: glkvm cloud access address (leave empty to auto-detect at startup)

   ⚠️ **Note:** All configuration should be done in the `.env` file.
    You don’t need to modify `docker-compose.yml`, templates, or scripts directly.

3. **Start the services**

   ```bash
   docker-compose up -d
   ```

   If you modify `.env` or template files, make sure to apply the updates:

   ```bash
   docker-compose down && docker-compose up -d
   ```

4. **Platform Access**

   Once the installation is complete, access the platform via: 

    ```bash
    https://<your_server_public_ip>
    ```
