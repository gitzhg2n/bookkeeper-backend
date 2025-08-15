package jobs

import (
	"context"
	"time"
	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
)

// LowBalanceJob checks for accounts below threshold and creates notifications
func LowBalanceJob(ctx context.Context, accountStore *db.AccountStore, userSettingsStore *db.UserSettingsStore, notificationStore *db.NotificationStore, userID uint) error {
	accounts, err := accountStore.ListByUser(userID)
	if err != nil {
		return err
	}
	settings, _ := userSettingsStore.GetByUserID(userID)
	threshold := int64(10000) // $100 default
	if settings != nil && settings.LowBalanceThreshold > 0 {
		threshold = settings.LowBalanceThreshold
	}
	for _, acc := range accounts {
		if acc.BalanceCents < threshold {
			msg := "Account '" + acc.Name + "' balance low: $" + formatCents(acc.BalanceCents)
			n := &models.Notification{
				UserID:  userID,
				Type:    models.NotificationType("low_balance"),
				Message: msg,
				Read:    false,
				CreatedAt: time.Now(),
			}
			notificationStore.CreateNotification(ctx, n)
		}
	}
	return nil
}

func formatCents(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}
