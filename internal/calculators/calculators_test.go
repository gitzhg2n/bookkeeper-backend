package calculators

import (
	"testing"
)

func TestMortgageMonthlyPayment(t *testing.T) {
	pay, err := MortgageMonthlyPayment(200000, 3.5, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pay <= 0 {
		t.Fatalf("expected positive payment, got %v", pay)
	}
}

func TestDebtPayoffMonthsAndInterest(t *testing.T) {
	months, interest, err := DebtPayoffMonthsAndInterest(5000, 10, 500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if months <= 0 || interest < 0 {
		t.Fatalf("unexpected results months=%d interest=%v", months, interest)
	}
}

func TestInvestmentFutureValue(t *testing.T) {
	fv, contrib, interest, err := InvestmentFutureValue(1000, 100, 5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fv <= 0 || contrib <= 0 || interest < 0 {
		t.Fatalf("unexpected results fv=%v contrib=%v interest=%v", fv, contrib, interest)
	}
}

func TestRentVsBuyComparison(t *testing.T) {
	netOwning, totalRenting, netBenefit, rec, err := RentVsBuyComparison(300000, 60000, 4, 30, 1.2, 1200, 1500, 2, 3000, 6, 1500, 300, 2, 10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if netOwning == 0 || totalRenting == 0 || rec == "" {
		t.Fatalf("unexpected results netOwning=%v totalRenting=%v rec=%v netBenefit=%v", netOwning, totalRenting, rec, netBenefit)
	}
}

func TestTaxEstimate(t *testing.T) {
	eff, tax, err := TaxEstimate(100000, 12000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eff <= 0 || tax <= 0 {
		t.Fatalf("unexpected results eff=%v tax=%v", eff, tax)
	}
}
