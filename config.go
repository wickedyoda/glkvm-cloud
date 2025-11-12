/*
 * MIT License
 *
 * Copyright (c) 2019 Jianhui Zhao <zhaojh329@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/kylelemons/go-gypsy/yaml"
    "github.com/urfave/cli/v3"
)

type Config struct {
    AddrDev              string
    AddrUser             string
    AddrHttpProxy        string
    HttpProxyRedirURL    string
    HttpProxyRedirDomain string
    Token                string
    DevHookUrl           string
    UserHookUrl          string
    LocalAuth            bool
    Password             string
    AllowOrigins         bool
    PprofAddr            string
    SslCert              string
    SslKey               string
    WebrtcIP             string
    WebrtcPort           string
    WebrtcUsername       string
    WebrtcPassword       string
    // LDAP配置 (LDAP Configuration)
    LdapEnabled       bool
    LdapServer        string
    LdapPort          int
    LdapUseTLS        bool
    LdapBindDN        string
    LdapBindPassword  string
    LdapBaseDN        string
    LdapUserFilter    string
    LdapAllowedGroups string
    LdapAllowedUsers  string

    // 通用OIDC Provider（支持任何标准OIDC提供商）
    OIDCEnabled             bool
    OIDCGenericClientID     string
    OIDCGenericClientSecret string
    OIDCGenericAuthURL      string
    OIDCGenericTokenURL     string
    OIDCGenericRedirectURL  string
    OIDCGenericScopes       []string
    OIDCGenericAllowedUsers []string
}

// docker mode fixed path for reading certificate
const (
    SslCert = "/home/certificate/glkvm_cer"
    SslKey  = "/home/certificate/glkvm_key"
)

func (cfg *Config) Parse(c *cli.Command) error {
    conf := c.String("conf")
    if conf != "" {
        err := parseYamlCfg(cfg, conf)
        if err != nil {
            return err
        }
    }

    getFlagOpt(c, "addr-dev", &cfg.AddrDev)
    getFlagOpt(c, "addr-user", &cfg.AddrUser)
    getFlagOpt(c, "addr-http-proxy", &cfg.AddrHttpProxy)
    getFlagOpt(c, "http-proxy-redir-url", &cfg.HttpProxyRedirURL)
    getFlagOpt(c, "http-proxy-redir-domain", &cfg.HttpProxyRedirDomain)
    getFlagOpt(c, "dev-hook-url", &cfg.DevHookUrl)
    getFlagOpt(c, "user-hook-url", &cfg.UserHookUrl)
    getFlagOpt(c, "local-auth", &cfg.LocalAuth)
    getFlagOpt(c, "token", &cfg.Token)
    getFlagOpt(c, "password", &cfg.Password)
    getFlagOpt(c, "allow-origins", &cfg.AllowOrigins)
    getFlagOpt(c, "pprof", &cfg.PprofAddr)

    cfg.SslCert = SslCert
    cfg.SslKey = SslKey

    getFlagOpt(c, "webrtc-ip", &cfg.WebrtcIP)
    getFlagOpt(c, "webrtc-port", &cfg.WebrtcPort)
    getFlagOpt(c, "webrtc-username", &cfg.WebrtcUsername)
    getFlagOpt(c, "webrtc-password", &cfg.WebrtcPassword)

    // LDAP配置标志 (LDAP Configuration flags)
    getFlagOpt(c, "ldap-enabled", &cfg.LdapEnabled)
    getFlagOpt(c, "ldap-server", &cfg.LdapServer)
    getFlagOpt(c, "ldap-port", &cfg.LdapPort)
    getFlagOpt(c, "ldap-use-tls", &cfg.LdapUseTLS)
    getFlagOpt(c, "ldap-bind-dn", &cfg.LdapBindDN)
    getFlagOpt(c, "ldap-bind-password", &cfg.LdapBindPassword)
    getFlagOpt(c, "ldap-base-dn", &cfg.LdapBaseDN)
    getFlagOpt(c, "ldap-user-filter", &cfg.LdapUserFilter)
    getFlagOpt(c, "ldap-allowed-groups", &cfg.LdapAllowedGroups)
    getFlagOpt(c, "ldap-allowed-users", &cfg.LdapAllowedUsers)

    // 通用OIDC Provider（支持任何标准OIDC提供商）
    getFlagOpt(c, "oidc-enabled", &cfg.OIDCEnabled)
    getFlagOpt(c, "oidc-generic-client-id", &cfg.OIDCGenericClientID)
    getFlagOpt(c, "oidc-generic-client-secret", &cfg.OIDCGenericClientSecret)
    getFlagOpt(c, "oidc-generic-auth-url", &cfg.OIDCGenericAuthURL)
    getFlagOpt(c, "oidc-generic-token-url", &cfg.OIDCGenericTokenURL)
    getFlagOpt(c, "oidc-generic-redirect-url", &cfg.OIDCGenericRedirectURL)
    getFlagOpt(c, "oidc-generic-scopes", &cfg.OIDCGenericScopes)
    getFlagOpt(c, "oidc-generic-allowed-users", &cfg.OIDCGenericAllowedUsers)

    return nil
}

func getConfigOpt(yamlCfg *yaml.File, name string, opt any) {
    val, err := yamlCfg.Get(name)
    if err != nil {
        return
    }

    switch opt := opt.(type) {
    case *string:
        *opt = val
    case *int:
        *opt, _ = strconv.Atoi(val)
    case *bool:
        *opt, _ = strconv.ParseBool(val)
    }
}

func parseYamlCfg(cfg *Config, conf string) error {
    yamlCfg, err := yaml.ReadFile(conf)
    if err != nil {
        return fmt.Errorf(`read config file: %s`, err.Error())
    }

    getConfigOpt(yamlCfg, "addr-dev", &cfg.AddrDev)
    getConfigOpt(yamlCfg, "addr-user", &cfg.AddrUser)
    getConfigOpt(yamlCfg, "addr-http-proxy", &cfg.AddrHttpProxy)
    getConfigOpt(yamlCfg, "http-proxy-redir-url", &cfg.HttpProxyRedirURL)
    getConfigOpt(yamlCfg, "http-proxy-redir-domain", &cfg.HttpProxyRedirDomain)

    getConfigOpt(yamlCfg, "token", &cfg.Token)
    getConfigOpt(yamlCfg, "dev-hook-url", &cfg.DevHookUrl)
    getConfigOpt(yamlCfg, "user-hook-url", &cfg.UserHookUrl)
    getConfigOpt(yamlCfg, "local-auth", &cfg.LocalAuth)
    getConfigOpt(yamlCfg, "password", &cfg.Password)
    getConfigOpt(yamlCfg, "allow-origins", &cfg.AllowOrigins)

    getConfigOpt(yamlCfg, "webrtc-ip", &cfg.WebrtcIP)
    getConfigOpt(yamlCfg, "webrtc-port", &cfg.WebrtcPort)
    getConfigOpt(yamlCfg, "webrtc-username", &cfg.WebrtcUsername)
    getConfigOpt(yamlCfg, "webrtc-password", &cfg.WebrtcPassword)

    // LDAP配置 (LDAP Configuration)
    getConfigOpt(yamlCfg, "ldap-enabled", &cfg.LdapEnabled)
    getConfigOpt(yamlCfg, "ldap-server", &cfg.LdapServer)
    getConfigOpt(yamlCfg, "ldap-port", &cfg.LdapPort)
    getConfigOpt(yamlCfg, "ldap-use-tls", &cfg.LdapUseTLS)
    getConfigOpt(yamlCfg, "ldap-bind-dn", &cfg.LdapBindDN)
    // 注意：为了避免特殊字符解析问题和提高安全性，ldap-bind-password故意不从YAML读取
    // Note: ldap-bind-password is intentionally not read from YAML to avoid special character parsing issues and for security.
    // It's always read directly from the LDAP_BIND_PASSWORD environment variable below
    getConfigOpt(yamlCfg, "ldap-base-dn", &cfg.LdapBaseDN)
    getConfigOpt(yamlCfg, "ldap-user-filter", &cfg.LdapUserFilter)
    getConfigOpt(yamlCfg, "ldap-allowed-groups", &cfg.LdapAllowedGroups)
    getConfigOpt(yamlCfg, "ldap-allowed-users", &cfg.LdapAllowedUsers)

    // LDAP密码始终从环境变量读取以避免YAML特殊字符解析问题和提高安全性
    // LDAP password is always read from environment variable to avoid YAML special character parsing issues and for security.
    if envPassword := os.Getenv("LDAP_BIND_PASSWORD"); envPassword != "" {
        cfg.LdapBindPassword = envPassword
    }

    // ===== OIDC 配置 (generic OIDC provider) =====
    // 开关与基本端点
    getConfigOpt(yamlCfg, "oidc-enabled", &cfg.OIDCEnabled)
    getConfigOpt(yamlCfg, "oidc-generic-client-id", &cfg.OIDCGenericClientID)
    getConfigOpt(yamlCfg, "oidc-generic-client-secret", &cfg.OIDCGenericClientSecret)

    getConfigOpt(yamlCfg, "oidc-generic-auth-url", &cfg.OIDCGenericAuthURL)
    getConfigOpt(yamlCfg, "oidc-generic-token-url", &cfg.OIDCGenericTokenURL)
    getConfigOpt(yamlCfg, "oidc-generic-redirect-url", &cfg.OIDCGenericRedirectURL)

    if s, err := yamlCfg.Get("oidc-generic-scopes"); err == nil && strings.TrimSpace(s) != "" {
        cfg.OIDCGenericScopes = splitScopes(s)
    }
    // 默认 scopes（当启用 OIDC 且仍未设置时）
    if cfg.OIDCEnabled && len(cfg.OIDCGenericScopes) == 0 {
        cfg.OIDCGenericScopes = []string{"openid", "profile", "email"}
    }

    if s, err := yamlCfg.Get("oidc-generic-allowed-users"); err == nil && strings.TrimSpace(s) != "" {
        cfg.OIDCGenericAllowedUsers = splitScopes(s)
    }

    return nil
}

func getFlagOpt(c *cli.Command, name string, opt any) {
    if !c.IsSet(name) {
        return
    }

    switch opt := opt.(type) {
    case *string:
        *opt = c.String(name)
    case *int:
        *opt = c.Int(name)
    case *bool:
        *opt = c.Bool(name)
    }
}

// splitScopes 把空格/逗号/换行分隔的 scopes 拆成去重、去空的列表
func splitScopes(s string) []string {
    // 先把逗号替换为空格，便于统一拆分
    s = strings.ReplaceAll(s, ",", " ")
    // 去掉 YAML 列表里可能的方括号等符号（宽容处理）
    s = strings.NewReplacer("[", " ", "]", " ").Replace(s)

    parts := strings.Fields(s)
    uniq := make([]string, 0, len(parts))
    seen := make(map[string]struct{}, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" || p == "-" { // 兼容误拷的 "-"
            continue
        }
        if _, ok := seen[p]; ok {
            continue
        }
        seen[p] = struct{}{}
        uniq = append(uniq, p)
    }
    return uniq
}
