package db

import (
	"context"
	"database/sql"
	"bookkeeper-backend/internal/models"
)

type AlertHistoryStore struct {
	DB *sql.DB
}

func (s *AlertHistoryStore) RecordAlertHistory(ctx context.Context, alertID, userID uint, details string) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO alert_history (alert_id, user_id, triggered_at, details) VALUES (?, ?, datetime('now'), ?)`,
		alertID, userID, details,
	)
	return err
}

func (s *AlertHistoryStore) ListAlertHistory(ctx context.Context, alertID uint) ([]models.AlertHistory, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, alert_id, user_id, triggered_at, details FROM alert_history WHERE alert_id = ? ORDER BY triggered_at DESC`, alertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.AlertHistory
	for rows.Next() {
		var h models.AlertHistory
		if err := rows.Scan(&h.ID, &h.AlertID, &h.UserID, &h.TriggeredAt, &h.Details); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
