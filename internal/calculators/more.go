package calculators

import (
	"errors"
	"math"
)

// LoanOffer represents a loan offer used in comparisons.
type LoanOffer struct {
	Principal    float64
	InterestRate float64 // annual percent
	TermYears    int
	Fees         float64 // upfront fees
}

// LoanComparisonResult contains per-offer results and ranking.
type LoanComparisonResult struct {
	MonthlyPayment float64
	TotalCost      float64 // total paid including fees
	TotalInterest  float64
}

// CompareLoans compares multiple loan offers returning results in same order.
func CompareLoans(offers []LoanOffer) ([]LoanComparisonResult, error) {
	if len(offers) == 0 {
		return nil, errors.New("no offers")
	}
	results := make([]LoanComparisonResult, len(offers))
	for i, o := range offers {
		if o.Principal <= 0 || o.TermYears <= 0 {
			return nil, errors.New("invalid offer")
		}
		monthly, err := MortgageMonthlyPayment(o.Principal, o.InterestRate, o.TermYears)
		if err != nil {
			return nil, err
		}
		months := o.TermYears * 12
		totalPaid := monthly * float64(months)
		totalCost := totalPaid + o.Fees
		totalInterest := totalPaid - o.Principal
		results[i] = LoanComparisonResult{MonthlyPayment: monthly, TotalCost: totalCost, TotalInterest: totalInterest}
	}
	return results, nil
}

// MortgageAffordability estimates max home price given debt-to-income constraints.
// Simplified: max monthly payment = (income * dtiRatio / 12) - otherMonthlyDebts
func MortgageAffordability(annualIncome, dtiRatio, downPayment, interestRate float64, loanTermYears int, otherMonthlyDebts float64) (float64, float64, error) {
	if annualIncome <= 0 || loanTermYears <= 0 {
		return 0, 0, errors.New("invalid inputs")
	}
	monthlyIncome := annualIncome / 12
	maxMonthly := monthlyIncome*(dtiRatio/100) - otherMonthlyDebts
	if maxMonthly <= 0 {
		return 0, 0, errors.New("no affordability with given params")
	}
	// invert mortgage payment formula to find loan principal
	monthlyRate := (interestRate / 100) / 12
	n := float64(loanTermYears * 12)
	var principal float64
	if monthlyRate == 0 {
		principal = maxMonthly * n
	} else {
		factor := math.Pow(1+monthlyRate, n)
		principal = (maxMonthly * (factor - 1)) / (monthlyRate * factor)
	}
	purchasePrice := principal + downPayment
	return purchasePrice, maxMonthly, nil
}

// CreditOptimization suggests allocation: choose snowball or avalanche based on inputs and returns recommendation.
func CreditOptimization(cards []CardInput, extraMonthly float64) (string, error) {
	if len(cards) == 0 {
		return "", errors.New("no cards")
	}
	// simple heuristic: if small number of cards and balances similar, avalanche; if many small balances -> snowball
	totalBal := 0.0
	minBal := math.Inf(1)
	maxBal := 0.0
	for _, c := range cards {
		totalBal += c.Balance
		if c.Balance < minBal {
			minBal = c.Balance
		}
		if c.Balance > maxBal {
			maxBal = c.Balance
		}
	}
	avg := totalBal / float64(len(cards))
	if maxBal > avg*2 {
		return "avalanche", nil
	}
	if minBal < avg*0.5 {
		return "snowball", nil
	}
	return "avalanche", nil
}

// CollegeSavings projects tuition inflation and required monthly contribution to reach target.
func CollegeSavings(currentBalance, annualTuition, tuitionInflation, annualReturn float64, years int) (float64, float64, error) {
	if years <= 0 {
		return 0, 0, errors.New("invalid years")
	}
	// project tuition at time of college start
	tuition := annualTuition * math.Pow(1+tuitionInflation/100, float64(years))
	// assume need for 4 years of tuition
	totalNeed := 0.0
	for i := 0; i < 4; i++ {
		totalNeed += tuition * math.Pow(1+tuitionInflation/100, float64(i))
	}
	// compute required monthly contribution to reach totalNeed given currentBalance and annualReturn
	rate := annualReturn / 100
	months := years * 12
	factor := math.Pow(1+rate/12, float64(months))
	var required float64
	if rate == 0 {
		required = (totalNeed - currentBalance) / float64(months)
	} else {
		required = (totalNeed - currentBalance*factor) * (rate/12) / (factor - 1)
	}
	return totalNeed, required, nil
}

// InvestmentFeeDrag computes impact of fees over time: returns difference in final balances.
func InvestmentFeeDrag(initial, annualContribution, grossReturn, feePercent float64, years int) (float64, float64, float64, error) {
	if years <= 0 {
		return 0, 0, 0, errors.New("invalid years")
	}
	rGross := grossReturn / 100
	fee := feePercent / 100
	fvGross := initial
	fvNet := initial
	for i := 0; i < years; i++ {
		fvGross = fvGross*(1+rGross) + annualContribution
		fvNet = fvNet*(1+(rGross-fee)) + annualContribution
	}
	return fvGross, fvNet, fvGross - fvNet, nil
}

// SafeWithdrawalRule applies a simple constant percentage withdrawal (e.g., 4%) to initial portfolio and projects remaining balance after years with return.
func SafeWithdrawal(initial, withdrawalRate, annualReturn float64, years int) (float64, error) {
	if years < 0 {
		return 0, errors.New("invalid years")
	}
	r := annualReturn / 100
	withdraw := initial * (withdrawalRate / 100)
	bal := initial
	for i := 0; i < years; i++ {
		bal = bal*(1+r) - withdraw
	}
	return bal, nil
}

// CD Ladder: returns slice of maturities and amounts given ladder size and total amount.
func CDLadder(totalAmount float64, rungs int, baseRate float64) ([]float64, error) {
	if totalAmount <= 0 || rungs <= 0 {
		return nil, errors.New("invalid inputs")
	}
	per := totalAmount / float64(rungs)
	maturities := make([]float64, rungs)
	for i := 0; i < rungs; i++ {
		rate := baseRate + float64(i)*0.1
		maturities[i] = per * math.Pow(1+rate/100, float64(i+1))
	}
	return maturities, nil
}

// Payroll simple estimator: net per period after flat rates
func PayrollEstimator(grossAnnual float64, payPeriods int, federalRate, stateRate float64, preTax float64) (float64, error) {
	if payPeriods <= 0 {
		return 0, errors.New("invalid pay periods")
	}
	annualTaxable := math.Max(0, grossAnnual-preTax)
	annualTax := annualTaxable*(federalRate/100) + annualTaxable*(stateRate/100)
	netAnnual := grossAnnual - annualTax - preTax
	return netAnnual / float64(payPeriods), nil
}

// Currency conversion using static rate (placeholder)
func ConvertCurrency(amount, rate float64) (float64, error) {
	if rate <= 0 {
		return 0, errors.New("invalid rate")
	}
	return amount * rate, nil
}
