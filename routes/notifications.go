package routes

import (
	"encoding/json"
	"net/http"

	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/middleware"
)

// NotificationHandler provides HTTP handlers for notifications
type NotificationHandler struct {
	Store *db.NotificationStore
}

// GET /notifications - list notifications for the authenticated user
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	notifications, err := h.Store.ListNotifications(r.Context(), int64(user.ID))
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notifications)
}

// POST /notifications/read - mark a notification as read
func (h *NotificationHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		NotificationID int64 `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Store.MarkNotificationRead(r.Context(), req.NotificationID, int64(user.ID)); err != nil {
		http.Error(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /notifications/read-all - mark all notifications as read
func (h *NotificationHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Store.MarkAllNotificationsRead(r.Context(), int64(user.ID)); err != nil {
		http.Error(w, "Failed to mark all as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
