package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"
)

// NotificationHandler provides HTTP handlers for notifications
type NotificationHandler struct {
	Store *db.NotificationStore
}

// GET /notifications - list notifications for the authenticated user
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	notifications, err := h.Store.ListNotifications(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notifications)
}

// POST /notifications/read - mark a notification as read
func (h *NotificationHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req struct {
		NotificationID int64 `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Store.MarkNotificationRead(r.Context(), req.NotificationID, userID); err != nil {
		http.Error(w, "Failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /notifications/read-all - mark all notifications as read
func (h *NotificationHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if err := h.Store.MarkAllNotificationsRead(r.Context(), userID); err != nil {
		http.Error(w, "Failed to mark all as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
