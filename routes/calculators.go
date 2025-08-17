package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"bookkeeper-backend/internal/calculators"
	"bookkeeper-backend/internal/models"
)

// MortgageCalculator handles mortgage payment calculations.
func MortgageCalculator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.MortgageCalculatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.LoanAmount <= 0 || req.InterestRate <= 0 || req.LoanTermYears <= 0 {
		writeJSONError(r, w, "loan amount, interest rate, and loan term must be positive", http.StatusBadRequest)
		return
	}

	monthlyPayment, err := calculators.MortgageMonthlyPayment(req.LoanAmount, req.InterestRate, req.LoanTermYears)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.MortgageCalculatorResponse{MonthlyPayment: monthlyPayment})
}

// DebtPayoffCalculator calculates months to pay off a debt and total interest paid.
func DebtPayoffCalculator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.DebtPayoffCalculatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.DebtAmount <= 0 || req.InterestRate < 0 || req.MonthlyPayment <= 0 {
		writeJSONError(r, w, "invalid inputs", http.StatusBadRequest)
		return
	}

	months, totalInterest, err := calculators.DebtPayoffMonthsAndInterest(req.DebtAmount, req.InterestRate, req.MonthlyPayment)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.DebtPayoffCalculatorResponse{MonthsToPayOff: months, TotalInterestPaid: totalInterest})
}

// InvestmentGrowthCalculator projects investment future value with monthly contributions.
func InvestmentGrowthCalculator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.InvestmentGrowthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.YearsToGrow <= 0 || req.InterestRate < 0 || (req.InitialPrincipal <= 0 && req.MonthlyContribution <= 0) {
		writeJSONError(r, w, "invalid inputs", http.StatusBadRequest)
		return
	}
	fv, totalContrib, totalInterest, err := calculators.InvestmentFutureValue(req.InitialPrincipal, req.MonthlyContribution, req.InterestRate, req.YearsToGrow)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.InvestmentGrowthResponse{FutureValue: fv, TotalContributions: totalContrib, TotalInterestEarned: totalInterest})
}

// RentVsBuyCalculator compares total cost of renting vs owning over time.
func RentVsBuyCalculator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.RentVsBuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.ComparisonYears <= 0 {
		writeJSONError(r, w, "invalid comparison years", http.StatusBadRequest)
		return
	}
	netOwningCost, totalRenting, netBenefit, rec, err := calculators.RentVsBuyComparison(req.HomePrice, req.DownPayment, req.InterestRate, req.LoanTermYears, req.PropertyTaxRate, req.HomeInsurance, req.MaintenanceCosts, req.AppreciationRate, req.ClosingCosts, req.SellingCostsRate, req.MonthlyRent, req.RentersInsurance, req.AnnualRentIncrease, req.ComparisonYears, req.InvestmentReturnRate)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.RentVsBuyResponse{TotalCostOfOwning: netOwningCost, TotalCostOfRenting: totalRenting, NetBenefitOfOwning: netBenefit, Recommendation: rec})
}

// TaxEstimatorCalculator provides a very simplified tax estimate (US-like progressive bands placeholder).
func TaxEstimatorCalculator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.TaxEstimatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	eff, tax, err := calculators.TaxEstimate(req.GrossIncome, req.Deductions)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.TaxEstimatorResponse{EffectiveTaxRate: eff, TotalTaxAmount: tax})
}

// AmortizationScheduleHandler returns full amortization schedule.
func AmortizationScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.AmortizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	start := req.StartDate
	if start.IsZero() {
		start = time.Now()
	}
	schedule, totalInterest, err := calculators.AmortizationSchedule(req.LoanAmount, req.InterestRate, req.LoanTermYears, start, req.ExtraMonthly)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	// Convert schedule rows to generic maps for JSON serialization
	out := make([]map[string]interface{}, 0, len(schedule))
	for _, row := range schedule {
		out = append(out, map[string]interface{}{
			"payment_date": row.PaymentDate,
			"payment": row.Payment,
			"principal": row.Principal,
			"interest": row.Interest,
			"remaining_balance": row.RemainingBalance,
		})
	}
	writeJSONSuccess(r, w, "ok", models.AmortizationResponse{Schedule: out, TotalInterest: totalInterest})
}

// RefinanceBreakevenHandler computes breakeven and savings for refinance.
func RefinanceBreakevenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.RefinanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	newMonthly, monthlySavings, breakeven, totalSavings, err := calculators.RefinanceBreakeven(req.CurrentBalance, req.CurrentRate, req.CurrentRemainingYears, req.NewRate, req.NewTermYears, req.ClosingCosts, req.YearsToEvaluate)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.RefinanceResponse{NewMonthly: newMonthly, MonthlySavings: monthlySavings, BreakevenMonths: breakeven, TotalSavings: totalSavings})
}

// APR/APY conversion handlers
func APRToAPYHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.APRAPYRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	apy, err := calculators.APRToAPY(req.Rate, req.CompoundsPerYear)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.APRAPYResponse{Converted: apy})
}

func APYToAPRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.APRAPYRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	apr, err := calculators.APYToAPR(req.Rate, req.CompoundsPerYear)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.APRAPYResponse{Converted: apr})
}

// RetirementProjectionHandler
func RetirementProjectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.RetirementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	nominal, real, yearly, err := calculators.RetirementProjection(req.CurrentBalance, req.AnnualSalary, req.ContributionRate, req.EmployerMatchRate, req.EmployerMatchCapPercent, req.AnnualReturn, req.Inflation, req.Years)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.RetirementResponse{NominalFutureValue: nominal, RealFutureValue: real, YearlyProjections: yearly})
}

// SavingsGoalHandler
func SavingsGoalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.SavingsGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	months, required, err := calculators.SavingsGoal(req.GoalAmount, req.CurrentSavings, req.MonthlyContribution, req.AnnualReturn)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.SavingsGoalResponse{MonthsToGoal: months, RequiredMonthly: required})
}

// Credit card payoff handler
func CreditPayoffHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.CreditPayoffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	cards := make([]calculators.CardInput, 0, len(req.Cards))
	for _, c := range req.Cards {
		cards = append(cards, calculators.CardInput{Balance: c.Balance, APR: c.APR, MinPayment: c.MinPayment})
	}
	res, err := calculators.CreditCardPayoff(cards, req.ExtraMonthly, req.Method)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.CreditPayoffResponse{MonthsToPayOff: res.MonthsToPayOff, TotalInterest: res.TotalInterest, Order: res.Order})
}

// TakeHomeHandler
func TakeHomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.TakeHomeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	netPer, err := calculators.TakeHomePay(req.GrossAnnual, req.PayPeriodsPerYear, req.FederalTaxRate, req.StateTaxRate, req.PreTaxDeductions)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.TakeHomeResponse{NetPerPeriod: netPer})
}

// InflationAdjustHandler
func InflationAdjustHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.InflationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	real, err := calculators.InflationAdjusted(req.Nominal, req.InflationRate, req.Years)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.InflationResponse{Real: real})
}

// NetWorthHandler
func NetWorthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.NetWorthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	net, alloc, err := calculators.NetWorthSnapshot(req.Assets, req.Liabilities)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.NetWorthResponse{NetWorth: net, Allocation: alloc})
}

// Loan comparison
func LoanComparisonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.LoanComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	offers := make([]calculators.LoanOffer, 0, len(req.Offers))
	for _, o := range req.Offers {
		offers = append(offers, calculators.LoanOffer{Principal: o.Principal, InterestRate: o.InterestRate, TermYears: o.TermYears, Fees: o.Fees})
	}
	res, err := calculators.CompareLoans(offers)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	out := make([]map[string]float64, 0, len(res))
	for _, r := range res {
		out = append(out, map[string]float64{"monthly_payment": r.MonthlyPayment, "total_cost": r.TotalCost, "total_interest": r.TotalInterest})
	}
	writeJSONSuccess(r, w, "ok", models.LoanComparisonResponse{Results: out})
}

// Mortgage affordability
func AffordabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.AffordabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	price, maxMonthly, err := calculators.MortgageAffordability(req.AnnualIncome, req.DTIRatio, req.DownPayment, req.InterestRate, req.LoanTermYears, req.OtherMonthlyDebts)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.AffordabilityResponse{PurchasePrice: price, MaxMonthly: maxMonthly})
}

// Credit optimization
func CreditOptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.CreditOptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	cards := make([]calculators.CardInput, 0, len(req.Cards))
	for _, c := range req.Cards {
		cards = append(cards, calculators.CardInput{Balance: c.Balance, APR: c.APR, MinPayment: c.MinPayment})
	}
	rec, err := calculators.CreditOptimization(cards, req.ExtraMonthly)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.CreditOptResponse{Recommendation: rec})
}

// College savings
func CollegeSavingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.CollegeSavingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	totalNeed, required, err := calculators.CollegeSavings(req.CurrentBalance, req.AnnualTuition, req.TuitionInflation, req.AnnualReturn, req.YearsUntilCollege)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.CollegeSavingsResponse{TotalNeed: totalNeed, RequiredMonthly: required})
}

// Fee drag
func FeeDragHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.FeeDragRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	gross, net, diff, err := calculators.InvestmentFeeDrag(req.Initial, req.AnnualContribution, req.GrossReturn, req.FeePercent, req.Years)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.FeeDragResponse{Gross: gross, Net: net, Difference: diff})
}

// Safe withdrawal
func SafeWithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.SafeWithdrawalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	rem, err := calculators.SafeWithdrawal(req.Initial, req.WithdrawalRate, req.AnnualReturn, req.Years)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.SafeWithdrawalResponse{Remaining: rem})
}

// CD ladder
func CDLadderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.CDLadderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	maturities, err := calculators.CDLadder(req.TotalAmount, req.Rungs, req.BaseRate)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.CDLadderResponse{Maturities: maturities})
}

// Payroll
func PayrollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.PayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	net, err := calculators.PayrollEstimator(req.GrossAnnual, req.PayPeriods, req.FederalRate, req.StateRate, req.PreTax)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.PayrollResponse{NetPerPeriod: net})
}

// Currency convert
func ConvertCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req models.ConvertCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	v, err := calculators.ConvertCurrency(req.Amount, req.Rate)
	if err != nil {
		writeJSONError(r, w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSONSuccess(r, w, "ok", models.ConvertCurrencyResponse{Converted: v})
}

