package db

import (
	"context"
	"database/sql"

	"bookkeeper-backend/internal/models"
)

type InvestmentAlertStore struct {
	DB *sql.DB
}

func (s *InvestmentAlertStore) ListActiveAlerts(ctx context.Context) ([]models.InvestmentAlert, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, user_id, rule, active, created_at, updated_at FROM investment_alerts WHERE active = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.InvestmentAlert
	for rows.Next() {
		var a models.InvestmentAlert
		if err := rows.Scan(&a.ID, &a.UserID, &a.Rule, &a.Active, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}
