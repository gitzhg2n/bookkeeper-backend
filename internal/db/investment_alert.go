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

// ListByUser returns all alerts for a specific user
func (s *InvestmentAlertStore) ListByUser(ctx context.Context, userID uint) ([]models.InvestmentAlert, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, user_id, asset_symbol, rule, alert_type, direction, threshold, active, created_at, updated_at FROM investment_alerts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.InvestmentAlert
	for rows.Next() {
		var a models.InvestmentAlert
		if err := rows.Scan(&a.ID, &a.UserID, &a.AssetSymbol, &a.Rule, &a.AlertType, &a.Direction, &a.Threshold, &a.Active, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

// Create inserts a new investment alert
func (s *InvestmentAlertStore) Create(ctx context.Context, a *models.InvestmentAlert) error {
	query := `INSERT INTO investment_alerts (user_id, asset_symbol, rule, alert_type, direction, threshold, active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	now := a.CreatedAt
	if now.IsZero() {
		now = now
	}
	return s.DB.QueryRowContext(ctx, query, a.UserID, a.AssetSymbol, a.Rule, a.AlertType, a.Direction, a.Threshold, a.Active, a.CreatedAt, a.UpdatedAt).Scan(&a.ID)
}

// Update modifies an existing alert (only by owner)
func (s *InvestmentAlertStore) Update(ctx context.Context, a *models.InvestmentAlert) error {
	query := `UPDATE investment_alerts SET asset_symbol=$1, rule=$2, alert_type=$3, direction=$4, threshold=$5, active=$6, updated_at=$7 WHERE id=$8 AND user_id=$9`
	_, err := s.DB.ExecContext(ctx, query, a.AssetSymbol, a.Rule, a.AlertType, a.Direction, a.Threshold, a.Active, a.UpdatedAt, a.ID, a.UserID)
	return err
}

// Delete removes an alert owned by user
func (s *InvestmentAlertStore) Delete(ctx context.Context, userID, id uint) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM investment_alerts WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
