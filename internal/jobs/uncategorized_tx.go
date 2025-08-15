package jobs

import (
	"context"
	"time"
	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
)

// UncategorizedTxJob finds transactions older than 3 days without a category and notifies the user
func UncategorizedTxJob(ctx context.Context, txStore *db.TransactionStore, notificationStore *db.NotificationStore, userID uint) error {
	cutoff := time.Now().AddDate(0, 0, -3)
	txs, err := txStore.ListUncategorizedBefore(userID, cutoff)
	if err != nil {
		return err
	}
	for _, tx := range txs {
		msg := "Uncategorized transaction: $" + formatCents(tx.AmountCents) + " on " + tx.OccurredAt.Format("2006-01-02")
		n := &models.Notification{
			UserID:  userID,
			Type:    models.NotificationType("uncategorized_tx"),
			Message: msg,
			Read:    false,
			CreatedAt: time.Now(),
		}
		notificationStore.CreateNotification(ctx, n)
	}
	return nil
}

func formatCents(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}
