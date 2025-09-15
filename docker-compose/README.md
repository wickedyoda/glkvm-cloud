# Quick Start

This is a simple guide for deploying **glkvm-cloud** using Docker Compose.

1. **Copy the environment file**

   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` and update the required parameters**

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

