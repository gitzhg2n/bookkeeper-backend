package jobs

import (
	"context"
	"time"

	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
)

// MarketDataProvider defines methods for fetching real-time market data, news, and sentiment
type MarketDataProvider interface {
	GetCurrentPrice(symbol string) (float64, error)
	GetPercentChange(symbol string, duration string) (float64, error)
	GetNewsSentiment(symbol string) (float64, error)
}

// ExampleStubMarketData is a stub implementation for development/testing
// Replace with real API integration in production
type ExampleStubMarketData struct{}

func (e *ExampleStubMarketData) GetCurrentPrice(symbol string) (float64, error) {
	return 100.0, nil
}
func (e *ExampleStubMarketData) GetPercentChange(symbol string, duration string) (float64, error) {
	return 5.0, nil
}
func (e *ExampleStubMarketData) GetNewsSentiment(symbol string) (float64, error) {
	return 0.0, nil // Neutral
}

// EvaluateInvestmentAlertsJob checks all user-configured investment alerts and triggers notifications if conditions are met.
func EvaluateInvestmentAlertsJob(ctx context.Context, dbStore *db.Store) error {
	alerts, err := dbStore.InvestmentAlertStore.ListActiveAlerts(ctx)
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		user, err := dbStore.UserStore.GetUserByID(ctx, int64(alert.UserID))
		if err != nil {
			continue
		}
		if user.Plan != "premium" && user.Plan != "selfhost" {
			continue // Only premium and selfhost users get advanced alerts
		}

		// Evaluate alert condition (advanced logic)
		triggered, details := EvaluateAlertCondition(ctx, dbStore, alert)
		if triggered {
			n := &models.Notification{
				UserID:    int64(alert.UserID),
				Type:      models.NotificationTypeInvestmentAlert,
				Title:     "Investment Alert Triggered",
				Message:   details,
				CreatedAt: time.Now(),
				Read:      false,
			}
			dbStore.NotificationStore.CreateNotification(ctx, n)
			// Record alert history
			if dbStore.AlertHistoryStore != nil {
				dbStore.AlertHistoryStore.RecordAlertHistory(ctx, alert.ID, alert.UserID, details)
			}
		}
	}
	return nil
}

// EvaluateAlertCondition evaluates the alert's rule and returns (triggered, details).
// Now supports compound/chained/time-based/custom rule logic for premium/selfhost users.
func EvaluateAlertCondition(ctx context.Context, dbStore *db.Store, alert models.InvestmentAlert) (bool, string) {
	// If compound condition is present, evaluate recursively
	if alert.Compound != nil {
		return evaluateCompoundCondition(ctx, dbStore, *alert.Compound)
	}
	// Example: Trigger if price crosses threshold (pseudo-code)
	// In real implementation, fetch current price from market data provider
	currentPrice := fetchCurrentPrice(alert.AssetSymbol)
	switch alert.AlertType {
	case "price":
		if (alert.Direction == "up" && currentPrice >= alert.Threshold) ||
			(alert.Direction == "down" && currentPrice <= alert.Threshold) {
			return true, "Price alert triggered for " + alert.AssetSymbol
		}
	case "percent_change":
		// Example: check percent change over period (pseudo-code)
		percentChange := fetchPercentChange(alert.AssetSymbol)
		if (alert.Direction == "up" && percentChange >= alert.Threshold) ||
			(alert.Direction == "down" && percentChange <= alert.Threshold) {
			return true, "Percent change alert triggered for " + alert.AssetSymbol
		}
	case "value":
		// Example: check portfolio value (pseudo-code)
		value := fetchPortfolioValue(alert.UserID, alert.AssetSymbol)
		if (alert.Direction == "up" && value >= alert.Threshold) ||
			(alert.Direction == "down" && value <= alert.Threshold) {
			return true, "Value alert triggered for " + alert.AssetSymbol
		}
	}
	// If time window is present, evaluate change over that period
	if alert.TimeWindow != nil {
		duration, err := time.ParseDuration(alert.TimeWindow.Duration)
		if err == nil {
			// Example: fetch price/percent change/value over the time window
			switch alert.AlertType {
			case "price":
				oldPrice := fetchHistoricalPrice(alert.AssetSymbol, duration)
				currentPrice := fetchCurrentPrice(alert.AssetSymbol)
				change := currentPrice - oldPrice
				if (alert.Direction == "up" && change >= alert.Threshold) ||
					(alert.Direction == "down" && change <= -alert.Threshold) {
					return true, "Price changed by threshold in time window for " + alert.AssetSymbol
				}
			case "percent_change":
				percentChange := fetchHistoricalPercentChange(alert.AssetSymbol, duration)
				if (alert.Direction == "up" && percentChange >= alert.Threshold) ||
					(alert.Direction == "down" && percentChange <= -alert.Threshold) {
					return true, "Percent change met threshold in time window for " + alert.AssetSymbol
				}
			case "value":
				oldValue := fetchHistoricalPortfolioValue(alert.UserID, alert.AssetSymbol, duration)
				currentValue := fetchPortfolioValue(alert.UserID, alert.AssetSymbol)
				change := currentValue - oldValue
				if (alert.Direction == "up" && change >= alert.Threshold) ||
					(alert.Direction == "down" && change <= -alert.Threshold) {
					return true, "Value changed by threshold in time window for " + alert.AssetSymbol
				}
			}
		}
	}
	// Check cooldown: do not trigger if a notification for this alert was sent within the cooldown period
	if alert.CooldownMinutes > 0 {
		recent, err := wasAlertRecentlyTriggered(ctx, dbStore, alert, alert.CooldownMinutes)
		if err == nil && recent {
			return false, ""
		}
	}
	// If custom rule is present, evaluate it
	if alert.CustomRule != "" {
		triggered, details, err := evaluateCustomRule(ctx, dbStore, alert)
		if err == nil && triggered {
			return true, details
		}
	}
	// TODO: Add compound/chained/time-based/custom rule logic here
	return false, ""
}

// evaluateCompoundCondition recursively evaluates AND/OR logic for compound alert conditions
func evaluateCompoundCondition(ctx context.Context, dbStore *db.Store, compound models.CompoundAlertCondition) (bool, string) {
	if len(compound.Conditions) == 0 {
		return false, ""
	}
	var triggered bool
	var details []string
	if compound.Operator == "AND" {
		triggered = true
		for _, cond := range compound.Conditions {
			singleAlert := models.InvestmentAlert{
				UserID:      cond.UserID,
				AssetSymbol: cond.AssetSymbol,
				AlertType:   cond.AlertType,
				Direction:   cond.Direction,
				Threshold:   cond.Threshold,
			}
			t, d := EvaluateAlertCondition(ctx, dbStore, singleAlert)
			if !t {
				triggered = false
			}
			if d != "" {
				details = append(details, d)
			}
		}
	} else if compound.Operator == "OR" {
		for _, cond := range compound.Conditions {
			singleAlert := models.InvestmentAlert{
				UserID:      cond.UserID,
				AssetSymbol: cond.AssetSymbol,
				AlertType:   cond.AlertType,
				Direction:   cond.Direction,
				Threshold:   cond.Threshold,
			}
			t, d := EvaluateAlertCondition(ctx, dbStore, singleAlert)
			if t {
				triggered = true
				if d != "" {
					details = append(details, d)
				}
				break
			}
		}
	}
	if triggered {
		return true, "Compound alert triggered: " + joinDetails(details)
	}
	return false, ""
}

func joinDetails(details []string) string {
	if len(details) == 0 {
		return ""
	}
	result := details[0]
	for i := 1; i < len(details); i++ {
		result += "; " + details[i]
	}
	return result
}

// Dummy functions for illustration (replace with real data sources)
func fetchCurrentPrice(symbol string) float64 { return 100.0 }
func fetchPercentChange(symbol string) float64 { return 5.0 }
func fetchPortfolioValue(userID uint, symbol string) float64 { return 1000.0 }
func fetchHistoricalPrice(symbol string, duration time.Duration) float64 { return 90.0 }
func fetchHistoricalPercentChange(symbol string, duration time.Duration) float64 { return 3.0 }
func fetchHistoricalPortfolioValue(userID uint, symbol string, duration time.Duration) float64 { return 950.0 }

// Helper to check if alert was triggered within cooldown period
func wasAlertRecentlyTriggered(ctx context.Context, dbStore *db.Store, alert models.InvestmentAlert, cooldownMinutes int) (bool, error) {
	since := time.Now().Add(-time.Duration(cooldownMinutes) * time.Minute)
	notifications, err := dbStore.NotificationStore.ListNotifications(ctx, int64(alert.UserID))
	if err != nil {
		return false, err
	}
	for _, n := range notifications {
		if n.Type == models.NotificationTypeInvestmentAlert && n.Title == "Investment Alert Triggered" && n.CreatedAt.After(since) {
			// Optionally, match on alert.AssetSymbol or other fields for more granularity
			return true, nil
		}
	}
	return false, nil
}

// Evaluate custom rule using a simple expression evaluator (pseudo-code)
func evaluateCustomRule(ctx context.Context, dbStore *db.Store, alert models.InvestmentAlert) (bool, string, error) {
	// For demonstration, support expressions like: "price > 100 && percent_change < -5"
	// In production, use a safe expression evaluator (e.g., github.com/Knetic/govaluate)
	// Here, we just check for a hardcoded example
	if alert.CustomRule == "price > 100 && percent_change < -5" {
		price := fetchCurrentPrice(alert.AssetSymbol)
		percentChange := fetchPercentChange(alert.AssetSymbol)
		if price > 100 && percentChange < -5 {
			return true, "Custom rule triggered: price > 100 && percent_change < -5", nil
		}
		return false, "", nil
	}
	// Add more parsing/evaluation as needed
	return false, "", nil
}
