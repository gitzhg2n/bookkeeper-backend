package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"
	"bookkeeper-backend/security"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Email = sanitizeString(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSONError(w, "email and password required", http.StatusBadRequest)
		return
	}

	if err := security.ValidatePasswordStrength(req.Password, h.cfg.AllowInsecurePassword); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare argon params
	argonParams := security.ArgonParams{
		MemoryKiB:   h.cfg.PasswordMemoryKiB,
		Time:        h.cfg.PasswordTime,
		Parallelism: h.cfg.PasswordParallelism,
		SaltLength:  h.cfg.PasswordSaltLength,
		KeyLength:   h.cfg.PasswordKeyLength,
	}

	argonSalt, err := security.RandomBytes(int(argonParams.SaltLength))
	if err != nil {
		writeJSONError(w, "internal error (salt gen)", http.StatusInternalServerError)
		return
	}
	passwordKey := security.DeriveKey(req.Password, argonSalt, argonParams)

	// KDF to produce KEK (reuse derived passwordKey as KEK for now – later we can HKDF)
	kek := passwordKey

	// Wrap fresh DEK
	_, encryptedDEK, err := security.WrapDEK(kek)
	if err != nil {
		writeJSONError(w, "internal error (wrap)", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:            req.Email,
		PasswordHash:     passwordKey, // For MVP we are storing raw Argon derived key; later we can store PHC formatted string
		EncryptedDEK:     encryptedDEK.Ciphertext,
		DEKNonce:         encryptedDEK.Nonce,
		ArgonMemoryKiB:   argonParams.MemoryKiB,
		ArgonTime:        argonParams.Time,
		ArgonParallelism: argonParams.Parallelism,
		ArgonSalt:        argonSalt,
		ArgonKeyLength:   argonParams.KeyLength,
		KDFVersion:       h.cfg.EncryptionKeyVersion,
	}

	if err := h.db.Create(user).Error; err != nil {
		writeJSONError(w, "user create failed (maybe duplicate email)", http.StatusConflict)
		return
	}

	accessToken, refreshToken, exp, err := h.issueTokens(user)
	if err != nil {
		writeJSONError(w, "token issue failed", http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, "registered", authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    exp,
		UserID:       user.ID,
		Email:        user.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}
	req.Email = sanitizeString(req.Email)
	if req.Email == "" || req.Password == "" {
		writeJSONError(w, "email and password required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		writeJSONError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Re-derive password key
	argonParams := security.ArgonParams{
		MemoryKiB:   user.ArgonMemoryKiB,
		Time:        user.ArgonTime,
		Parallelism: user.ArgonParallelism,
		SaltLength:  uint32(len(user.ArgonSalt)),
		KeyLength:   user.ArgonKeyLength,
	}
	key := security.DeriveKey(req.Password, user.ArgonSalt, argonParams)
	if !security.ConstantTimeCompare(key, user.PasswordHash) {
		writeJSONError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, exp, err := h.issueTokens(&user)
	if err != nil {
		writeJSONError(w, "token issue failed", http.StatusInternalServerError)
		return
	}

	writeJSONSuccess(w, "authenticated", authResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    exp,
		UserID:       user.ID,
		Email:        user.Email,
	})
}

func (h *AuthHandler) issueTokens(u *models.User) (accessToken, refreshToken string, expires time.Time, err error) {
	now := time.Now()
	accessExp := now.Add(h.cfg.AccessTokenTTL)
	refreshExp := now.Add(h.cfg.RefreshTokenTTL)

	accessClaims := middleware.Claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshClaims := middleware.Claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
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
	return at, rt, accessExp, nil
}

// Placeholder for refresh endpoint logic (coming soon)
var ErrNotImplemented = errors.New("not implemented")