package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/fanjindong/go-cache"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/sessions"
    "github.com/rs/zerolog/log"
    "io"
    "math/rand"
    "net/http"
    "net/url"
    "rttys/utils"
    "strings"
    "time"
)

var (
    // Session store for OAuth flow
    oauthStore *sessions.CookieStore
)

// Register OIDC routes
func RegisterOIDCRoutes(r *gin.Engine, cfg *Config) {
    if !cfg.OIDCEnabled {
        return
    }

    // Initialize session store
    sessionSecret := generateRandomString(32)
    oauthStore = sessions.NewCookieStore([]byte(sessionSecret))

    // Configure session options
    oauthStore.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   300, // 5 minutes, enough to complete the OAuth flow
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode, // Important: allow cookies on cross-site navigation (OIDC redirect)
        Domain:   "",                   // Empty means current host/domain will be used
    }

    // OIDC auth routes (public, no existing auth required)
    r.GET("/auth/oidc/login", oidcLoginHandler(cfg))
    r.GET("/auth/oidc/callback", oidcCallbackHandler(cfg))
}

// Start OIDC login
func oidcLoginHandler(cfg *Config) gin.HandlerFunc {
    return func(c *gin.Context) {

        // Generate state and nonce
        nonce := generateRandomString(32)
        state := generateRandomString(32)

        log.Info().Msgf("OIDC login initiated: nonce=%s, state=%s...", nonce[:10], state[:10])

        // Save to session
        session, _ := oauthStore.Get(c.Request, "oidc-session")
        session.Values["state"] = state
        session.Values["nonce"] = nonce

        if err := session.Save(c.Request, c.Writer); err != nil {
            log.Error().Err(err).Msg("Failed to save OIDC session")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
            return
        }

        // Build authorization URL
        params := url.Values{}
        params.Add("client_id", cfg.OIDCGenericClientID)
        params.Add("redirect_uri", cfg.OIDCGenericRedirectURL)
        params.Add("response_type", "code")
        params.Add("scope", strings.Join(cfg.OIDCGenericScopes, " "))
        params.Add("state", state)
        params.Add("nonce", nonce)

        authURL := cfg.OIDCGenericAuthURL + "?" + params.Encode()
        log.Info().Msgf("OIDC login redirect=%s", authURL[:100]+"...")
        c.Redirect(http.StatusFound, authURL)
    }
}

// Handle OIDC callback
func oidcCallbackHandler(cfg *Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get session
        session, err := oauthStore.Get(c.Request, "oidc-session")
        if err != nil {
            log.Error().Err(err).Msg("Failed to get OIDC session")
            c.Redirect(http.StatusFound, "/?error=session_error")
            return
        }

        // Validate state
        state := c.Query("state")
        savedState, ok := session.Values["state"].(string)
        if !ok || state != savedState {
            log.Warn().Msg("OIDC state mismatch")
            c.Redirect(http.StatusFound, "/?error=invalid_state")
            return
        }

        // Check OAuth error
        if errorMsg := c.Query("error"); errorMsg != "" {
            errorDesc := c.Query("error_description")
            log.Warn().Msgf("OIDC error: %s - %s", errorMsg, errorDesc)
            c.Redirect(http.StatusFound, "/?error="+errorMsg)
            return
        }

        // Read authorization code
        code := c.Query("code")
        if code == "" {
            c.Redirect(http.StatusFound, "/?error=no_code")
            return
        }

        // Exchange authorization code for tokens
        tokens, err := exchangeCodeForTokens(cfg, code)
        if err != nil {
            log.Error().Err(err).Msg("Failed to exchange code for tokens")
            c.Redirect(http.StatusFound, "/?error=token_exchange_failed")
            return
        }

        // Extract user info
        var userEmail string
        var userName string

        // Standard OIDC – parse from ID token
        idToken, ok := tokens["id_token"].(string)
        if !ok {
            log.Error().Msg("No ID token in response")
            c.Redirect(http.StatusFound, "/?error=no_id_token")
            return
        }

        claims, err := parseIDToken(idToken)
        if err != nil {
            log.Error().Err(err).Msg("Failed to parse ID token")
            c.Redirect(http.StatusFound, "/?error=invalid_token")
            return
        }

        // Validate nonce
        if savedNonce, ok := session.Values["nonce"].(string); ok {
            if claims["nonce"] != savedNonce {
                log.Warn().Msg("OIDC nonce mismatch")
                c.Redirect(http.StatusFound, "/?error=invalid_nonce")
                return
            }
        }

        userEmail, _ = claims["email"].(string)
        userName, _ = claims["name"].(string)
        if userName == "" {
            userName = userEmail
        }

        if err != nil || userEmail == "" {
            log.Error().Err(err).Msg("Failed to get user info")
            c.Redirect(http.StatusFound, "/?error=user_info_failed")
            return
        }

        log.Info().Msgf("OIDC login successful:  email=%s", userEmail)

        // ====== OIDC whitelist enforcement ======
        if len(cfg.OIDCGenericAllowedUsers) > 0 && !isEmailAllowed(cfg.OIDCGenericAllowedUsers, userEmail) {
            log.Warn().Msgf("OIDC user not in whitelist: email=%s", userEmail)
            // Differentiate authorization failure types on frontend via error code
            c.Redirect(http.StatusFound, "/?error=authorization")
            return
        }

        // Create application session
        sid := utils.GenUniqueID()
        httpSessions.Set(sid, gin.H{
            "email": userEmail,
            "name":  userName,
            "oidc":  true,
        }, cache.WithEx(httpSessionExpire))

        c.SetCookie("sid", sid, 0, "", "", cfg.SslCert != "", true)

        // Clean up OAuth session
        session.Options.MaxAge = -1
        session.Save(c.Request, c.Writer)

        // Redirect to home page
        c.Redirect(http.StatusFound, "/")
    }
}

// Exchange authorization code for tokens
func exchangeCodeForTokens(cfg *Config, code string) (map[string]interface{}, error) {
    data := url.Values{}
    data.Set("code", code)
    data.Set("client_id", cfg.OIDCGenericClientID)
    data.Set("client_secret", cfg.OIDCGenericClientSecret)
    data.Set("redirect_uri", cfg.OIDCGenericRedirectURL)
    data.Set("grant_type", "authorization_code")

    req, err := http.NewRequest("POST", cfg.OIDCGenericTokenURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Accept", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    log.Info().Msgf("OIDC exchange code for tokens successful:  %s", body)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("token request failed: %s", string(body))
    }

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result, nil
}

// Parse ID Token
func parseIDToken(idToken string) (map[string]interface{}, error) {
    parts := strings.Split(idToken, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid token format")
    }

    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, err
    }

    var claims map[string]interface{}
    if err := json.Unmarshal(payload, &claims); err != nil {
        return nil, err
    }

    return claims, nil
}

// Generate a random string with given length
func generateRandomString(length int) string {
    b := make([]byte, length)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)[:length]
}

// Email whitelist check rules:
// - Exact match (case-insensitive)
// - Simple domain match: entries starting with "@example.com" or "*@example.com"
//   mean "allow all users under this domain"
func isEmailAllowed(allowed []string, email string) bool {
    if len(allowed) == 0 {
        return true
    }
    e := strings.ToLower(strings.TrimSpace(email))
    for _, raw := range allowed {
        a := strings.ToLower(strings.TrimSpace(raw))
        if a == "" {
            continue
        }
        if strings.HasPrefix(a, "*@") || strings.HasPrefix(a, "@") {
            // Domain match
            dom := strings.TrimPrefix(a, "*@")
            dom = strings.TrimPrefix(dom, "@")
            if strings.HasSuffix(e, "@"+dom) {
                return true
            }
            continue
        }
        // Exact match
        if e == a {
            return true
        }
    }
    return false
}
