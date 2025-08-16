package calculators

import (
	"testing"
	"time"
)

func TestAmortizationSchedule(t *testing.T) {
	schedule, totalInterest, err := AmortizationSchedule(100000, 4, 30, time.Now(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schedule) == 0 || totalInterest <= 0 {
		t.Fatalf("unexpected results schedule=%d interest=%v", len(schedule), totalInterest)
	}
}

func TestRefinanceBreakeven(t *testing.T) {
	newMonthly, monthlySavings, breakeven, totalSavings, err := RefinanceBreakeven(200000, 4, 25, 3.5, 30, 3000, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newMonthly <= 0 || monthlySavings <= 0 {
		t.Fatalf("unexpected results %v %v", newMonthly, monthlySavings)
	}
	if breakeven <= 0 && totalSavings >= 0 {
		t.Fatalf("breakeven weird: %v %v", breakeven, totalSavings)
	}
}

func TestAPRAPY(t *testing.T) {
	apy, err := APRToAPY(5, 12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apy <= 0 {
		t.Fatalf("bad apy %v", apy)
	}
	apr, err := APYToAPR(apy, 12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apr <= 0 {
		t.Fatalf("bad apr %v", apr)
	}
}

func TestRetirementProjection(t *testing.T) {
	nominal, real, yearsVec, err := RetirementProjection(10000, 80000, 5, 50, 3, 6, 2, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nominal <= 0 || real <= 0 || len(yearsVec) != 30 {
		t.Fatalf("unexpected results nominal=%v real=%v len=%d", nominal, real, len(yearsVec))
	}
}

func TestSavingsGoal(t *testing.T) {
	months, req, err := SavingsGoal(10000, 1000, 100, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if months <= 0 && req <= 0 {
		t.Fatalf("unexpected results months=%v req=%v", months, req)
	}
}

func TestCreditCardPayoff(t *testing.T) {
	cards := []CardInput{{Balance: 5000, APR: 18, MinPayment: 150}, {Balance: 2000, APR: 22, MinPayment: 50}}
	res, err := CreditCardPayoff(cards, 100, "avalanche")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MonthsToPayOff <= 0 || res.TotalInterest <= 0 {
		t.Fatalf("unexpected results %v", res)
	}
}

func TestTakeHomePay(t *testing.T) {
	net, err := TakeHomePay(80000, 24, 12, 5, 2000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net <= 0 {
		t.Fatalf("bad net %v", net)
	}
}

func TestInflationAdjusted(t *testing.T) {
	real, err := InflationAdjusted(200000, 2, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if real <= 0 {
		t.Fatalf("bad real %v", real)
	}
}

func TestNetWorthSnapshot(t *testing.T) {
	assets := map[string]float64{"cash": 10000, "investments": 50000}
	liabs := map[string]float64{"mortgage": 150000}
	net, alloc, err := NetWorthSnapshot(assets, liabs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net >= 0 {
		t.Fatalf("expected negative net in this test, got %v", net)
	}
	if len(alloc) == 0 {
		t.Fatalf("expected allocation map")
	}
}
