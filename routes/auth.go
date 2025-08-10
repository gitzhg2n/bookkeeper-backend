package routes

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"
	"bookkeeper-backend/internal/security"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg    *config.Config
	db     *gorm.DB
	logger *slog.Logger
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, logger: logger}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
	Email        string    `json:"email"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Email = sanitizeString(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSONError(r, w, "email and password required", http.StatusBadRequest)
		return
	}

	if err := security.ValidatePasswordStrength(req.Password, h.cfg.AllowInsecurePassword); err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}

	argonParams := security.ArgonParams{
		MemoryKiB:   h.cfg.PasswordMemoryKiB,
		Time:        h.cfg.PasswordTime,
		Parallelism: h.cfg.PasswordParallelism,
		SaltLength:  h.cfg.PasswordSaltLength,
		KeyLength:   h.cfg.PasswordKeyLength,
	}

	argonSalt, err := security.RandomBytes(int(argonParams.SaltLength))
	if err != nil {
		writeJSONError(r, w, "internal error", http.StatusInternalServerError)
		return
	}
	passwordKey := security.DeriveKey(req.Password, argonSalt, argonParams)
 copilot/fix-184f7982-e511-4e6f-9dc2-305d1c6b4c15
	kek := security.DeriveKEK(passwordKey, "bookkeeper:dek:v1")

copilot/fix-bf106389-f58d-4461-b471-056cdc30d4c5
	kek := security.DeriveKEK(passwordKey, "bookkeeper:dek:v1")

	kek, err := security.DeriveKEK(passwordKey, "bookkeeper:dek:v1")
	if err != nil {
		writeJSONError(r, w, "internal error", http.StatusInternalServerError)
		return
	}
 main
 main
	_, encDEK, err := security.WrapDEK(kek)
	if err != nil {
		writeJSONError(r, w, "internal error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:            req.Email,
		PasswordHash:     passwordKey,
		EncryptedDEK:     encDEK.Ciphertext,
		DEKNonce:         encDEK.Nonce,
		ArgonMemoryKiB:   argonParams.MemoryKiB,
		ArgonTime:        argonParams.Time,
		ArgonParallelism: argonParams.Parallelism,
		ArgonSalt:        argonSalt,
		ArgonKeyLength:   argonParams.KeyLength,
		KDFVersion:       h.cfg.EncryptionKeyVersion,
	}

	if err := h.db.Create(user).Error; err != nil {
		writeJSONError(r, w, "user create failed", http.StatusConflict)
		return
	}

	at, rt, exp, err := h.issueTokens(user)
	if err != nil {
		writeJSONError(r, w, "token issue failed", http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(r, w, "registered", authResponse{
		AccessToken:  at,
		RefreshToken: rt,
		ExpiresAt:    exp,
		UserID:       user.ID,
		Email:        user.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Email = sanitizeString(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSONError(r, w, "email and password required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		writeJSONError(r, w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	argonParams := security.ArgonParams{
		MemoryKiB:   user.ArgonMemoryKiB,
		Time:        user.ArgonTime,
		Parallelism: user.ArgonParallelism,
		SaltLength:  uint32(len(user.ArgonSalt)),
		KeyLength:   user.ArgonKeyLength,
	}
	key := security.DeriveKey(req.Password, user.ArgonSalt, argonParams)
	if !security.ConstantTimeCompare(key, user.PasswordHash) {
		writeJSONError(r, w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	at, rt, exp, err := h.issueTokens(&user)
	if err != nil {
		writeJSONError(r, w, "token issue failed", http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(r, w, "authenticated", authResponse{
		AccessToken:  at,
		RefreshToken: rt,
		ExpiresAt:    exp,
		UserID:       user.ID,
		Email:        user.Email,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		writeJSONError(r, w, "refresh_token required", http.StatusBadRequest)
		return
	}
	token, claims, err := h.parseToken(req.RefreshToken)
	if err != nil || !token.Valid {
		writeJSONError(r, w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	var rt models.RefreshToken
	if err := h.db.Where("id = ? AND user_id = ?", claims.ID, claims.UserID).First(&rt).Error; err != nil {
		writeJSONError(r, w, "refresh invalid", http.StatusUnauthorized)
		return
	}
	if rt.RevokedAt != nil || time.Now().Unix() > rt.ExpiresAt {
		writeJSONError(r, w, "refresh expired or revoked", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		writeJSONError(r, w, "user not found", http.StatusUnauthorized)
		return
	}

	at, newRT, exp, err := h.issueTokens(&user)
	if err != nil {
		writeJSONError(r, w, "token issue failed", http.StatusInternalServerError)
		return
	}
	nowUnix := time.Now().Unix()
	if err := h.db.Model(&rt).Updates(map[string]any{
		"revoked_at":     &nowUnix,
		"replaced_by_id": claims.RegisteredClaims.ID,
	}).Error; err != nil {
		h.logger.Warn("failed revoke refresh", "error", err)
	}

	writeJSONSuccess(r, w, "refreshed", authResponse{
		AccessToken:  at,
		RefreshToken: newRT,
		ExpiresAt:    exp,
		UserID:       user.ID,
		Email:        user.Email,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req logoutRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.RefreshToken == "" {
		writeJSONError(r, w, "refresh_token required", http.StatusBadRequest)
		return
	}
	_, claims, err := h.parseToken(req.RefreshToken)
	if err != nil {
		writeJSONError(r, w, "invalid token", http.StatusUnauthorized)
		return
	}
	now := time.Now().Unix()
	h.db.Model(&models.RefreshToken{}).Where("id = ? AND user_id = ?", claims.ID, claims.UserID).
		Updates(map[string]any{"revoked_at": &now})
	writeJSONSuccess(r, w, "logged out", nil)
}

func (h *AuthHandler) issueTokens(u *models.User) (accessToken, refreshToken string, expires time.Time, err error) {
	now := time.Now()
	accessExp := now.Add(h.cfg.AccessTokenTTL)
	refreshExp := now.Add(h.cfg.RefreshTokenTTL)

	jti := uuid.NewString()
	refreshJTI := uuid.NewString()

	accessClaims := middleware.Claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshClaims := middleware.Claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJTI,
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(h.cfg.JWTSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(h.cfg.JWTSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	rec := &models.RefreshToken{
		ID:        refreshJTI,
		UserID:    u.ID,
		ExpiresAt: refreshExp.Unix(),
	}
	if err := h.db.Create(rec).Error; err != nil {
		return "", "", time.Time{}, err
	}
	return at, rt, accessExp, nil
}

func (h *AuthHandler) parseToken(tokenStr string) (*jwt.Token, *middleware.Claims, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return h.cfg.JWTSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

var ErrNotImplemented = errors.New("not implemented")