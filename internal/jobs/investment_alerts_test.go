package jobs_test

import (
	"context"
	"testing"
	"time"

	"bookkeeper-backend/internal/jobs"
	"bookkeeper-backend/internal/models"
)

func TestEvaluateAlertCondition_SimplePrice(t *testing.T) {
	alert := models.InvestmentAlert{
		UserID:      1,
		AssetSymbol: "BTC",
		AlertType:   "price",
		Direction:   "up",
		Threshold:   50.0,
	}
	triggered, details := jobs.EvaluateAlertCondition(context.Background(), nil, alert)
	if !triggered {
		t.Errorf("Expected alert to trigger, got not triggered. Details: %s", details)
	}
}

func TestEvaluateAlertCondition_CompoundAND(t *testing.T) {
	compound := &models.CompoundAlertCondition{
		Operator: "AND",
		Conditions: []models.InvestmentAlertCondition{
			{AssetSymbol: "BTC", AlertType: "price", Direction: "up", Threshold: 50.0},
			{AssetSymbol: "ETH", AlertType: "price", Direction: "up", Threshold: 50.0},
		},
	}
	alert := models.InvestmentAlert{
		UserID:   1,
		Compound: compound,
	}
	triggered, _ := jobs.EvaluateAlertCondition(context.Background(), nil, alert)
	if !triggered {
		t.Error("Expected compound AND alert to trigger")
	}
}

func TestEvaluateAlertCondition_TimeWindow(t *testing.T) {
	alert := models.InvestmentAlert{
		UserID:      1,
		AssetSymbol: "BTC",
		AlertType:   "price",
		Direction:   "up",
		Threshold:   5.0,
		TimeWindow:  &models.TimeWindow{Duration: "24h"},
	}
	triggered, _ := jobs.EvaluateAlertCondition(context.Background(), nil, alert)
	if !triggered {
		t.Error("Expected time window alert to trigger")
	}
}

func TestEvaluateAlertCondition_CustomRule(t *testing.T) {
	alert := models.InvestmentAlert{
		UserID:      1,
		AssetSymbol: "BTC",
		CustomRule:  "price > 100 && percent_change < -5",
	}
	triggered, details := jobs.EvaluateAlertCondition(context.Background(), nil, alert)
	if !triggered {
		t.Errorf("Expected custom rule alert to trigger, got not triggered. Details: %s", details)
	}
}
