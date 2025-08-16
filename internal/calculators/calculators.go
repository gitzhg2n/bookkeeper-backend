package calculators

import (
	"errors"
	"math"
)

// MortgageMonthlyPayment returns the monthly payment for a loan.
// interestRate is annual percent (e.g., 3.5 for 3.5%)
func MortgageMonthlyPayment(loanAmount, interestRate float64, loanTermYears int) (float64, error) {
	if loanAmount <= 0 || loanTermYears <= 0 {
		return 0, errors.New("loan amount and term must be positive")
	}
	monthlyInterestRate := (interestRate / 100) / 12
	numberOfPayments := loanTermYears * 12
	if monthlyInterestRate == 0 {
		return loanAmount / float64(numberOfPayments), nil
	}
	p := loanAmount
	i := monthlyInterestRate
	n := float64(numberOfPayments)
	numerator := p * i * math.Pow(1+i, n)
	denominator := math.Pow(1+i, n) - 1
	if denominator == 0 {
		return 0, errors.New("calculation error")
	}
	return numerator / denominator, nil
}

// DebtPayoffMonthsAndInterest simulates month-by-month payoff and returns months and total interest.
func DebtPayoffMonthsAndInterest(debtAmount, interestRate, monthlyPayment float64) (int, float64, error) {
	if debtAmount <= 0 || monthlyPayment <= 0 || interestRate < 0 {
		return 0, 0, errors.New("invalid inputs")
	}
	monthlyRate := (interestRate / 100) / 12
	balance := debtAmount
	months := 0
	totalInterest := 0.0
	for balance > 0 && months < 1000*12 {
		interest := balance * monthlyRate
		principal := monthlyPayment - interest
		if principal <= 0 {
			return 0, 0, errors.New("monthly payment too low to ever pay off debt")
		}
		balance -= principal
		totalInterest += interest
		months++
	}
	if months >= 1000*12 {
		return 0, 0, errors.New("calculation did not converge")
	}
	return months, totalInterest, nil
}

// InvestmentFutureValue computes future value with monthly contributions.
func InvestmentFutureValue(initialPrincipal, monthlyContribution, interestRate float64, years int) (float64, float64, float64, error) {
	if years <= 0 || interestRate < -100 || (initialPrincipal <= 0 && monthlyContribution <= 0) {
		return 0, 0, 0, errors.New("invalid inputs")
	}
	rRate := interestRate / 100
	months := years * 12
	fv := initialPrincipal
	totalContrib := 0.0
	for m := 0; m < months; m++ {
		fv = fv*(1+rRate/12) + monthlyContribution
		totalContrib += monthlyContribution
	}
	totalInterest := fv - initialPrincipal - totalContrib
	return fv, totalContrib, totalInterest, nil
}

// RentVsBuyComparison computes simple costs and recommendation.
func RentVsBuyComparison(homePrice, downPayment, interestRate float64, loanTermYears int, propertyTaxRate, homeInsurance, maintenanceCosts, appreciationRate, closingCosts, sellingCostsRate float64, monthlyRent, rentersInsurance, annualRentIncrease float64, comparisonYears int, investmentReturnRate float64) (float64, float64, float64, string, error) {
	if comparisonYears <= 0 {
		return 0, 0, 0, "", errors.New("invalid comparison years")
	}
	loanAmount := homePrice - downPayment
	monthlyRate := (interestRate / 100) / 12
	n := loanTermYears * 12
	var monthlyMortgage float64
	if monthlyRate == 0 {
		monthlyMortgage = loanAmount / float64(n)
	} else {
		numerator := loanAmount * monthlyRate * math.Pow(1+monthlyRate, float64(n))
		denominator := math.Pow(1+monthlyRate, float64(n)) - 1
		if denominator == 0 {
			return 0, 0, 0, "", errors.New("calculation error")
		}
		monthlyMortgage = numerator / denominator
	}
	totalOwning := 0.0
	homeValue := homePrice
	for y := 0; y < comparisonYears; y++ {
		totalOwning += monthlyMortgage * 12
		totalOwning += homeValue * (propertyTaxRate / 100)
		totalOwning += homeInsurance + maintenanceCosts
		homeValue = homeValue * (1 + appreciationRate/100)
	}
	sellingCosts := homeValue * (sellingCostsRate / 100)
	saleProceeds := homeValue - sellingCosts
	netOwningCost := totalOwning + downPayment + closingCosts - saleProceeds

	totalRenting := 0.0
	mr := monthlyRent
	for y := 0; y < comparisonYears; y++ {
		totalRenting += mr * 12
		totalRenting += rentersInsurance
		mr = mr * (1 + annualRentIncrease/100)
	}
	investRate := investmentReturnRate / 100
	investFV := downPayment
	for y := 0; y < comparisonYears; y++ {
		investFV = investFV * (1 + investRate)
	}
	netBenefit := investFV - netOwningCost - totalRenting
	rec := "Rent"
	if netBenefit > 0 {
		rec = "Buy"
	}
	return netOwningCost, totalRenting, netBenefit, rec, nil
}

// TaxEstimate calculates a simplified progressive tax estimate.
func TaxEstimate(grossIncome, deductions float64) (float64, float64, error) {
	taxable := grossIncome - deductions
	if taxable < 0 {
		taxable = 0
	}
	tax := 0.0
	remaining := taxable
	bands := []struct{
		cap float64
		rate float64
	}{
		{10000, 0.10},
		{40000, 0.12},
		{85000, 0.22},
		{1e18, 0.24},
	}
	prevCap := 0.0
	for _, b := range bands {
		capAmount := math.Max(0, math.Min(remaining, b.cap - prevCap))
		tax += capAmount * b.rate
		remaining -= capAmount
		prevCap = b.cap
		if remaining <= 0 {
			break
		}
	}
	eff := 0.0
	if grossIncome > 0 {
		eff = tax / grossIncome
	}
	return eff, tax, nil
}
