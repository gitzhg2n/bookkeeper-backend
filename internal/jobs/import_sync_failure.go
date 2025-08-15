package jobs

import (
	"context"
	"time"
	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/internal/db"
)

// ImportSyncFailureJob is a stub for future import/sync failure notifications
func ImportSyncFailureJob(ctx context.Context, notificationStore *db.NotificationStore, userID uint, details string) error {
	n := &models.Notification{
		UserID:  userID,
		Type:    models.NotificationType("import_sync_failure"),
		Message: "Import/Sync failure: " + details,
		Read:    false,
		CreatedAt: time.Now(),
	}
	return notificationStore.CreateNotification(ctx, n)
}
