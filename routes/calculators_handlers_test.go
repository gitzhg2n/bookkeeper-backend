package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculatorHandlersSmoke(t *testing.T) {
	// build a minimal mux similar to BuildRouter but with handlers directly accessible
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/calculators/mortgage", MortgageCalculator)
	mux.HandleFunc("/v1/calculators/debt-payoff", DebtPayoffCalculator)
	mux.HandleFunc("/v1/calculators/investment-growth", InvestmentGrowthCalculator)
	mux.HandleFunc("/v1/calculators/rent-vs-buy", RentVsBuyCalculator)
	mux.HandleFunc("/v1/calculators/tax-estimator", TaxEstimatorCalculator)
	mux.HandleFunc("/v1/calculators/amortization", AmortizationScheduleHandler)
	mux.HandleFunc("/v1/calculators/refinance-breakeven", RefinanceBreakevenHandler)
	mux.HandleFunc("/v1/calculators/apr-to-apy", APRToAPYHandler)
	mux.HandleFunc("/v1/calculators/apy-to-apr", APYToAPRHandler)
	mux.HandleFunc("/v1/calculators/retirement-projection", RetirementProjectionHandler)
	mux.HandleFunc("/v1/calculators/savings-goal", SavingsGoalHandler)
	mux.HandleFunc("/v1/calculators/credit-payoff", CreditPayoffHandler)
	mux.HandleFunc("/v1/calculators/take-home", TakeHomeHandler)
	mux.HandleFunc("/v1/calculators/inflation-adjust", InflationAdjustHandler)
	mux.HandleFunc("/v1/calculators/net-worth", NetWorthHandler)
	mux.HandleFunc("/v1/calculators/loan-comparison", LoanComparisonHandler)
	mux.HandleFunc("/v1/calculators/affordability", AffordabilityHandler)
	mux.HandleFunc("/v1/calculators/credit-optimization", CreditOptHandler)
	mux.HandleFunc("/v1/calculators/college-savings", CollegeSavingsHandler)
	mux.HandleFunc("/v1/calculators/fee-drag", FeeDragHandler)
	mux.HandleFunc("/v1/calculators/safe-withdrawal", SafeWithdrawalHandler)
	mux.HandleFunc("/v1/calculators/cd-ladder", CDLadderHandler)
	mux.HandleFunc("/v1/calculators/payroll", PayrollHandler)
	mux.HandleFunc("/v1/calculators/convert-currency", ConvertCurrencyHandler)

	tests := []struct{
		path string
		body string
	}{
		{"/v1/calculators/mortgage", `{"loan_amount":200000,"interest_rate":3.5,"loan_term_years":30}`},
		{"/v1/calculators/debt-payoff", `{"debt_amount":5000,"interest_rate":10,"monthly_payment":500}`},
		{"/v1/calculators/investment-growth", `{"initial_principal":1000,"monthly_contribution":100,"interest_rate":5,"years_to_grow":10}`},
		{"/v1/calculators/rent-vs-buy", `{"home_price":300000,"down_payment":60000,"interest_rate":4,"loan_term_years":30,"property_tax_rate":1.2,"home_insurance":1200,"maintenance_costs":1500,"appreciation_rate":2,"closing_costs":3000,"selling_costs_rate":6,"monthly_rent":1500,"renters_insurance":300,"annual_rent_increase":2,"comparison_years":10,"investment_return_rate":5}`},
		{"/v1/calculators/tax-estimator", `{"filing_status":"single","gross_income":100000,"deductions":12000}`},
		{"/v1/calculators/amortization", `{"loan_amount":200000,"interest_rate":3.5,"loan_term_years":30}`},
		{"/v1/calculators/refinance-breakeven", `{"current_balance":200000,"current_rate":4,"current_remaining_years":25,"new_rate":3.5,"new_term_years":30,"closing_costs":3000,"years_to_evaluate":10}`},
		{"/v1/calculators/apr-to-apy", `{"rate":5,"compounds_per_year":12}`},
		{"/v1/calculators/apy-to-apr", `{"rate":5,"compounds_per_year":12}`},
		{"/v1/calculators/retirement-projection", `{"current_balance":10000,"annual_salary":80000,"contribution_rate":5,"employer_match_rate":50,"employer_match_cap_percent":3,"annual_return":6,"inflation":2,"years":30}`},
		{"/v1/calculators/savings-goal", `{"goal_amount":10000,"current_savings":1000,"monthly_contribution":100,"annual_return":5}`},
		{"/v1/calculators/credit-payoff", `{"cards":[{"balance":5000,"apr":18,"min_payment":150},{"balance":2000,"apr":22,"min_payment":50}],"extra_monthly":100,"method":"avalanche"}`},
		{"/v1/calculators/take-home", `{"gross_annual":80000,"pay_periods_per_year":24,"federal_tax_rate":12,"state_tax_rate":5,"pre_tax_deductions":2000}`},
		{"/v1/calculators/inflation-adjust", `{"nominal":200000,"inflation_rate":2,"years":10}`},
		{"/v1/calculators/net-worth", `{"assets":{"cash":10000,"investments":50000},"liabilities":{"mortgage":150000}}`},
		{"/v1/calculators/loan-comparison", `{"offers":[{"principal":200000,"interest_rate":3.5,"term_years":30,"fees":2000},{"principal":200000,"interest_rate":4.0,"term_years":30,"fees":0}]}`},
		{"/v1/calculators/affordability", `{"annual_income":90000,"dti_ratio":36,"down_payment":30000,"interest_rate":4,"loan_term_years":30,"other_monthly_debts":500}`},
		{"/v1/calculators/credit-optimization", `{"cards":[{"balance":5000,"apr":18,"min_payment":150},{"balance":1000,"apr":22,"min_payment":50}],"extra_monthly":100}`},
		{"/v1/calculators/college-savings", `{"current_balance":5000,"annual_tuition":20000,"tuition_inflation":3,"annual_return":5,"years_until_college":10}`},
		{"/v1/calculators/fee-drag", `{"initial":10000,"annual_contribution":5000,"gross_return":8,"fee_percent":1,"years":20}`},
		{"/v1/calculators/safe-withdrawal", `{"initial":100000,"withdrawal_rate":4,"annual_return":5,"years":30}`},
		{"/v1/calculators/cd-ladder", `{"total_amount":10000,"rungs":5,"base_rate":2}`},
		{"/v1/calculators/payroll", `{"gross_annual":80000,"pay_periods":24,"federal_rate":12,"state_rate":5,"pre_tax":2000}`},
		{"/v1/calculators/convert-currency", `{"amount":100,"rate":1.1}`},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("POST", tt.path, bytes.NewBufferString(tt.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s returned status %d", tt.path, resp.StatusCode)
		}
	}
}
