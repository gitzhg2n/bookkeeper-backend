package routes

import (
	"net/http"

	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	writeJSONSuccess(r, w, "ok", map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}