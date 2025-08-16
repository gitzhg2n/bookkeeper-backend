package calculators

import "testing"

func TestCompareLoans(t *testing.T) {
	offers := []LoanOffer{{Principal: 200000, InterestRate: 3.5, TermYears: 30, Fees: 2000}, {Principal: 200000, InterestRate: 4.0, TermYears: 30, Fees: 0}}
	res, err := CompareLoans(offers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results")
	}
}

func TestMortgageAffordability(t *testing.T) {
	price, maxMonthly, err := MortgageAffordability(90000, 36, 30000, 4, 30, 500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price <= 0 || maxMonthly <= 0 {
		t.Fatalf("unexpected results")
	}
}

func TestCreditOptimization(t *testing.T) {
	cards := []CardInput{{Balance: 5000, APR: 18, MinPayment: 150}, {Balance: 1000, APR: 22, MinPayment: 50}}
	rec, err := CreditOptimization(cards, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == "" {
		t.Fatalf("empty recommendation")
	}
}

func TestCollegeSavings(t *testing.T) {
	total, req, err := CollegeSavings(5000, 20000, 3, 5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total <= 0 || req < 0 {
		t.Fatalf("unexpected results")
	}
}

func TestInvestmentFeeDrag(t *testing.T) {
	gross, net, diff, err := InvestmentFeeDrag(10000, 5000, 8, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gross <= net || diff <= 0 {
		t.Fatalf("unexpected results")
	}
}

func TestSafeWithdrawal(t *testing.T) {
	bal, err := SafeWithdrawal(100000, 4, 5, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bal < 0 {
		t.Fatalf("unexpected negative balance: %v", bal)
	}
}

func TestCDLadder(t *testing.T) {
	m, err := CDLadder(10000, 5, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 5 {
		t.Fatalf("unexpected length")
	}
}

func TestPayrollEstimator(t *testing.T) {
	net, err := PayrollEstimator(80000, 24, 12, 5, 2000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if net <= 0 {
		t.Fatalf("bad net %v", net)
	}
}

func TestConvertCurrency(t *testing.T) {
	v, err := ConvertCurrency(100, 1.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 110 {
		t.Fatalf("unexpected conversion %v", v)
	}
}
