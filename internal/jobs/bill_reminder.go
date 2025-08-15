package jobs

import (
	"context"
	"time"
	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
)

// BillReminderJob checks for bills due in 3 days and creates notifications
func BillReminderJob(ctx context.Context, billStore *db.BillStore, notificationStore *db.NotificationStore, userID uint) error {
	bills, err := billStore.ListDueInDays(userID, 3)
	if err != nil {
		return err
	}
	for _, bill := range bills {
		msg := "Bill '" + bill.Name + "' is due soon: $" + formatCents(bill.AmountCents) + " on " + bill.NextDue.Format("2006-01-02")
		n := &models.Notification{
			UserID:  userID,
			Type:    models.NotificationType("bill"),
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
