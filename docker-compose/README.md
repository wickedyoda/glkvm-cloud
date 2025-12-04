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

   **LDAP Authentication (Optional):**

   - `LDAP_ENABLED`: set to `true` to enable LDAP authentication (default: `false`)
   - `LDAP_SERVER`: LDAP server hostname or IP address
   - `LDAP_PORT`: LDAP server port (default: `389`, for TLS use `636`)
   - `LDAP_USE_TLS`: set to `true` to enable TLS encryption (default: `false`)
   - `LDAP_BIND_DN`: service account distinguished name
   - `LDAP_BIND_PASSWORD`: service account password
   - `LDAP_BASE_DN`: search base for user queries
   - `LDAP_USER_FILTER`: LDAP query filter (default: `(uid=%s)`)
   - `LDAP_ALLOWED_GROUPS`: comma-separated list of authorized groups (optional)
   - `LDAP_ALLOWED_USERS`: comma-separated list of authorized users (optional)

   ⚠️ **Note:** All configuration should be done in the `.env` file.
    You don’t need to modify `docker-compose.yml`, templates, or scripts directly.

   **OIDC Authentication (Optional):**

   - `OIDC_ENABLED`: set to `true` to enable OIDC authentication (default: `false`)
   - `OIDC_ISSUER`: OIDC issuer URL provided by your identity provider
      (e.g. `https://accounts.google.com`, `https://your-tenant.auth0.com/`)
   - `OIDC_CLIENT_ID`: client ID issued by your OIDC provider
   - `OIDC_CLIENT_SECRET`: client secret issued by your OIDC provider
   - `OIDC_AUTH_URL`: authorization endpoint URL
   - `OIDC_TOKEN_URL`: token endpoint URL
   - `OIDC_REDIRECT_URL`: redirect (callback) URL registered in your OIDC provider
      Domain is user-defined, but the path must be fixed: `/auth/oidc/callback`
      Example: `https://your-domain.example.com/auth/oidc/callback`
   - `OIDC_SCOPES`: space-separated list of requested scopes (default: `"openid profile email"`)
   - `OIDC_ALLOWED_USERS`: comma-separated list of allowed emails or domains (optional)
      Example: `user@example.com,@example.com`
   - `OIDC_ALLOWED_SUBS`: comma-separated list of allowed OIDC subject (`sub`) IDs (optional)
   - `OIDC_ALLOWED_USERNAMES`: comma-separated list of allowed usernames (`preferred_username` or `name`) (optional)
   - `OIDC_ALLOWED_GROUPS`: comma-separated list of allowed OIDC groups (optional)
   
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
