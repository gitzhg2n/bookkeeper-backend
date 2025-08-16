package calculators

import (
	"errors"
	"math"
	"time"
)

// AmortizationRow represents a single payment row in an amortization schedule.
type AmortizationRow struct {
	PaymentDate    time.Time `json:"payment_date"`
	Payment        float64   `json:"payment"`
	Principal      float64   `json:"principal"`
	Interest       float64   `json:"interest"`
	RemainingBalance float64 `json:"remaining_balance"`
}

// AmortizationSchedule returns a slice of AmortizationRow given loan params.
func AmortizationSchedule(loanAmount, interestRate float64, loanTermYears int, startDate time.Time, extraMonthly float64) ([]AmortizationRow, float64, error) {
	if loanAmount <= 0 || loanTermYears <= 0 {
		return nil, 0, errors.New("loan amount and term must be positive")
	}
	monthlyPayment, err := MortgageMonthlyPayment(loanAmount, interestRate, loanTermYears)
	if err != nil {
		return nil, 0, err
	}
	if extraMonthly < 0 {
		extraMonthly = 0
	}
	monthlyPayment += extraMonthly
	monthlyRate := (interestRate / 100) / 12
	balance := loanAmount
	schedule := make([]AmortizationRow, 0)
	paymentDate := startDate
	totalInterest := 0.0
	maxIter := loanTermYears * 12 * 2
	for balance > 0 && len(schedule) < maxIter {
		interest := balance * monthlyRate
		principal := monthlyPayment - interest
		if principal <= 0 {
			return nil, 0, errors.New("monthly payment too low to amortize")
		}
		if principal > balance {
			principal = balance
			monthlyPayment = principal + interest
		}
		balance -= principal
		totalInterest += interest
		schedule = append(schedule, AmortizationRow{
			PaymentDate:      paymentDate,
			Payment:          monthlyPayment,
			Principal:        principal,
			Interest:         interest,
			RemainingBalance: math.Max(0, balance),
		})
		paymentDate = paymentDate.AddDate(0, 1, 0)
	}
	if len(schedule) >= maxIter && balance > 0 {
		return schedule, totalInterest, errors.New("did not amortize within iteration limit")
	}
	return schedule, totalInterest, nil
}

// RefinanceBreakeven computes monthly payment change, monthly savings, cumulative savings timeline and breakeven month.
// returns newMonthlyPayment, monthlySavings, breakevenMonths, totalSavingsAtYears
func RefinanceBreakeven(currentBalance, currentRate float64, currentRemainingYears int, newRate float64, newTermYears int, closingCosts float64, yearsToEvaluate int) (float64, float64, int, float64, error) {
	if currentBalance <= 0 || currentRemainingYears <= 0 || newTermYears <= 0 {
		return 0, 0, 0, 0, errors.New("invalid inputs")
	}
	currentMonthly, err := MortgageMonthlyPayment(currentBalance, currentRate, currentRemainingYears)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	newMonthly, err := MortgageMonthlyPayment(currentBalance, newRate, newTermYears)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	monthlySavings := currentMonthly - newMonthly
	if yearsToEvaluate <= 0 {
		yearsToEvaluate = newTermYears
	}
	cumulative := 0.0
	breakevenMonth := -1
	for m := 1; m <= yearsToEvaluate*12; m++ {
		cumulative += monthlySavings
		if cumulative >= closingCosts && breakevenMonth == -1 {
			breakevenMonth = m
			// don't break; we want total savings at horizon
		}
	}
	totalSavings := cumulative - closingCosts
	return newMonthly, monthlySavings, breakevenMonth, totalSavings, nil
}

// APRToAPY converts an APR nominal rate to APY given compounding periods per year.
func APRToAPY(apr float64, compoundsPerYear int) (float64, error) {
	if compoundsPerYear <= 0 {
		return 0, errors.New("invalid compounding frequency")
	}
	r := apr / 100
	apy := math.Pow(1+r/float64(compoundsPerYear), float64(compoundsPerYear)) - 1
	return apy * 100, nil
}

// APYToAPR converts APY to nominal APR given compounding periods per year.
func APYToAPR(apy float64, compoundsPerYear int) (float64, error) {
	if compoundsPerYear <= 0 {
		return 0, errors.New("invalid compounding frequency")
	}
	r := apy / 100
	apr := (math.Pow(1+r, 1/float64(compoundsPerYear)) - 1) * float64(compoundsPerYear)
	return apr * 100, nil
}

// RetirementProjection projects retirement savings with employer match and inflation-adjusted final value.
// returns nominalFV, realFV, yearlyProjections (slice of balances)
func RetirementProjection(currentBalance, annualSalary, contributionRate, employerMatchRate, employerMatchCapPercent, annualReturn, inflation float64, years int) (float64, float64, []float64, error) {
	if years <= 0 {
		return 0, 0, nil, errors.New("years must be positive")
	}
	yearlyContrib := annualSalary * (contributionRate / 100)
	employerMatch := math.Min(yearlyContrib*(employerMatchRate/100), annualSalary*(employerMatchCapPercent/100))
	fv := currentBalance
	yearsVec := make([]float64, 0, years)
	rate := annualReturn / 100
	for y := 0; y < years; y++ {
		fv = fv*(1+rate) + yearlyContrib + employerMatch
		yearsVec = append(yearsVec, fv)
	}
	realFV := fv / math.Pow(1+(inflation/100), float64(years))
	return fv, realFV, yearsVec, nil
}

// SavingsGoal computes either time to reach goal or required monthly contribution to hit goal.
// If monthlyContribution > 0 returns monthsToGoal, else returns requiredMonthlyContribution.
func SavingsGoal(goalAmount, currentSavings, monthlyContribution, annualReturn float64) (int, float64, error) {
	if goalAmount <= 0 {
		return 0, 0, errors.New("goal must be positive")
	}
	rate := annualReturn / 100
	if monthlyContribution > 0 {
		months := 0
		balance := currentSavings
		for balance < goalAmount && months < 1000*12 {
			balance = balance*(1+rate/12) + monthlyContribution
			months++
		}
		if months >= 1000*12 {
			return 0, 0, errors.New("did not reach goal within iteration limit")
		}
		return months, 0, nil
	}
	// compute required monthly contribution using future value of annuity formula
	// FV = P*(1+r)^n + PMT*(((1+r)^n -1)/r)
	// solve for PMT
	// assume years = 1..100, pick years such that it's reasonable? Instead return required monthly for 1 year horizon as example
	years := 1
	n := years * 12
	factor := math.Pow(1+rate/12, float64(n))
	if rate == 0 {
		required := (goalAmount - currentSavings) / float64(n)
		return n, required, nil
	}
	required := (goalAmount - currentSavings*factor) * (rate/12) / (factor - 1)
	return n, required, nil
}

// CreditCardPayoffPlan computes payoff under snowball or avalanche methods.
// cards: slice of (balance, apr, minPayment)
type CardInput struct {
	Balance    float64
	APR        float64
	MinPayment float64
}

type PayoffResult struct {
	MonthsToPayOff int
	TotalInterest  float64
	Order          []int // indexes in order paid
}

func CreditCardPayoff(cards []CardInput, extraMonthly float64, method string) (PayoffResult, error) {
	if len(cards) == 0 {
		return PayoffResult{}, errors.New("no cards")
	}
	if extraMonthly < 0 {
		extraMonthly = 0
	}
	// copy balances
	balances := make([]float64, len(cards))
	for i, c := range cards {
		balances[i] = c.Balance
		if c.MinPayment <= 0 {
			return PayoffResult{}, errors.New("min payment must be positive")
		}
	}
	months := 0
	totalInterest := 0.0
	order := make([]int, 0, len(cards))
	paid := make([]bool, len(cards))
	// simple simulation up to 1000*12 months
	for sumBalances(balances) > 0 && months < 1000*12 {
		months++
		// determine target card
		var idx int
		if method == "avalanche" {
			// highest APR
			high := -1.0
			for i, c := range cards {
				if balances[i] > 0 && c.APR > high {
					high = c.APR
					idx = i
				}
			}
		} else {
			// snowball: smallest balance
			minBal := math.Inf(1)
			for i := range cards {
				if balances[i] > 0 && balances[i] < minBal {
					minBal = balances[i]
					idx = i
				}
			}
		}
		// compute payments
		for i := range cards {
			if balances[i] <= 0 {
				continue
			}
			interest := balances[i] * (cards[i].APR/100) / 12
			totalInterest += interest
			payment := cards[i].MinPayment
			if i == idx {
				payment += extraMonthly
			}
			principal := payment - interest
			if principal <= 0 {
				return PayoffResult{}, errors.New("payment too small to cover interest")
			}
			if principal > balances[i] {
				principal = balances[i]
			}
			balances[i] -= principal
			if balances[i] == 0 && !paid[i] {
				order = append(order, i)
				paid[i] = true
			}
		}
	}
	if months >= 1000*12 {
		return PayoffResult{}, errors.New("did not converge")
	}
	return PayoffResult{MonthsToPayOff: months, TotalInterest: totalInterest, Order: order}, nil
}

func sumBalances(b []float64) float64 {
	s := 0.0
	for _, v := range b {
		s += v
	}
	return s
}

// TakeHomePay estimates net pay per period for simple federal/state tax assumptions (very approximate).
func TakeHomePay(grossAnnual float64, payPeriodsPerYear int, federalTaxRate float64, stateTaxRate float64, preTaxDeductions float64) (float64, error) {
	if payPeriodsPerYear <= 0 {
		return 0, errors.New("invalid pay periods")
	}
	annualTaxable := math.Max(0, grossAnnual - preTaxDeductions)
	annualTax := annualTaxable*(federalTaxRate/100) + annualTaxable*(stateTaxRate/100)
	netAnnual := grossAnnual - annualTax - preTaxDeductions
	return netAnnual / float64(payPeriodsPerYear), nil
}

// InflationAdjustedProjection converts nominal future value to real value given inflation rate.
func InflationAdjusted(nominal float64, inflationRate float64, years int) (float64, error) {
	if years < 0 {
		return 0, errors.New("invalid years")
	}
	real := nominal / math.Pow(1+(inflationRate/100), float64(years))
	return real, nil
}

// NetWorthSnapshot computes net worth and allocation percentages.
func NetWorthSnapshot(assets map[string]float64, liabilities map[string]float64) (float64, map[string]float64, error) {
	totalAssets := 0.0
	for _, v := range assets {
		totalAssets += v
	}
	totalLiabilities := 0.0
	for _, v := range liabilities {
		totalLiabilities += v
	}
	net := totalAssets - totalLiabilities
	alloc := make(map[string]float64)
	if totalAssets > 0 {
		for k, v := range assets {
			alloc[k] = (v / totalAssets) * 100
		}
	}
	return net, alloc, nil
}
