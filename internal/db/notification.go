package db

import (
	"context"
	"database/sql"

	"bookkeeper-backend/internal/models"
)

// NotificationStore handles DB operations for notifications
type NotificationStore struct {
	DB *sql.DB
}

// CreateNotification inserts a new notification for a user
func (s *NotificationStore) CreateNotification(ctx context.Context, n *models.Notification) error {
	query := `INSERT INTO notifications (user_id, type, title, message, read, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return s.DB.QueryRowContext(ctx, query, n.UserID, n.Type, n.Title, n.Message, n.Read, n.CreatedAt).Scan(&n.ID)
}

// ListNotifications returns all notifications for a user (most recent first)
func (s *NotificationStore) ListNotifications(ctx context.Context, userID int64) ([]models.Notification, error) {
	query := `SELECT id, user_id, type, title, message, read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		n := models.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

// MarkNotificationRead sets a notification as read
func (s *NotificationStore) MarkNotificationRead(ctx context.Context, notificationID, userID int64) error {
	query := `UPDATE notifications SET read = TRUE WHERE id = $1 AND user_id = $2`
	_, err := s.DB.ExecContext(ctx, query, notificationID, userID)
	return err
}

// MarkAllNotificationsRead sets all notifications as read for a user
func (s *NotificationStore) MarkAllNotificationsRead(ctx context.Context, userID int64) error {
	query := `UPDATE notifications SET read = TRUE WHERE user_id = $1 AND read = FALSE`
	_, err := s.DB.ExecContext(ctx, query, userID)
	return err
}
