package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/FourneauxThibaut/CF-Back/internal/config"
	"github.com/gin-gonic/gin"
)

// AuthProxy proxies auth requests to Supabase. Credentials stay server-side.
type AuthProxy struct {
	baseURL string
	anonKey string
}

// NewAuthProxy creates an auth proxy with the given config.
func NewAuthProxy(cfg *config.Config) *AuthProxy {
	return &AuthProxy{
		baseURL: strings.TrimSuffix(cfg.SupabaseURL, "/") + "/auth/v1",
		anonKey: cfg.SupabaseAnonKey,
	}
}

func (p *AuthProxy) do(method, path string, body []byte, authToken string) (*http.Response, error) {
	req, err := http.NewRequest(method, p.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", p.anonKey)
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	return http.DefaultClient.Do(req)
}

// LoginReq is the request body for POST /auth/login.
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /auth/login. Proxies to Supabase token endpoint.
func (p *AuthProxy) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	body, _ := json.Marshal(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})
	resp, err := p.do("POST", "/token?grant_type=password", body, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	p.forwardResponse(c, resp)
}

// SignupReq is the request body for POST /auth/signup.
type SignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup handles POST /auth/signup. Proxies to Supabase signup endpoint.
func (p *AuthProxy) Signup(c *gin.Context) {
	var req SignupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	body, _ := json.Marshal(map[string]any{
		"email":    req.Email,
		"password": req.Password,
	})
	resp, err := p.do("POST", "/signup", body, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	p.forwardResponse(c, resp)
}

// Logout handles POST /auth/logout. Proxies to Supabase logout.
func (p *AuthProxy) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization required"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header"})
		return
	}
	resp, err := p.do("POST", "/logout", nil, token)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
		return
	}
	io.Copy(io.Discard, resp.Body)
	c.JSON(resp.StatusCode, gin.H{"error": "logout failed"})
}

// RefreshReq is the request body for POST /auth/refresh.
type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh handles POST /auth/refresh. Proxies to Supabase token refresh.
func (p *AuthProxy) Refresh(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}
	body, _ := json.Marshal(map[string]string{
		"refresh_token": req.RefreshToken,
		"grant_type":    "refresh_token",
	})
	resp, err := p.do("POST", "/token?grant_type=refresh_token", body, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()
	p.forwardResponse(c, resp)
}

func (p *AuthProxy) forwardResponse(c *gin.Context, resp *http.Response) {
	data, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/json"
	}
	c.Data(resp.StatusCode, ct, data)
}
