package models

// MortgageCalculatorRequest defines the structure for a mortgage calculation request.
type MortgageCalculatorRequest struct {
	LoanAmount    float64 `json:"loan_amount"     binding:"required"`
	InterestRate  float64 `json:"interest_rate"   binding:"required"`
	LoanTermYears int     `json:"loan_term_years" binding:"required"`
}

// MortgageCalculatorResponse defines the structure for a mortgage calculation response.
type MortgageCalculatorResponse struct {
	MonthlyPayment float64 `json:"monthly_payment"`
}

// DebtPayoffCalculatorRequest defines the structure for a debt payoff calculation request.
type DebtPayoffCalculatorRequest struct {
	DebtAmount     float64 `json:"debt_amount"      binding:"required"`
	InterestRate   float64 `json:"interest_rate"    binding:"required"`
	MonthlyPayment float64 `json:"monthly_payment"  binding:"required"`
}

// DebtPayoffCalculatorResponse defines the structure for a debt payoff calculation response.
type DebtPayoffCalculatorResponse struct {
	MonthsToPayOff    int     `json:"months_to_pay_off"`
	TotalInterestPaid float64 `json:"total_interest_paid"`
}

// InvestmentGrowthRequest defines the structure for an investment growth calculation request.
type InvestmentGrowthRequest struct {
	InitialPrincipal    float64 `json:"initial_principal"`
	MonthlyContribution float64 `json:"monthly_contribution" binding:"required"`
	InterestRate        float64 `json:"interest_rate"         binding:"required"`
	YearsToGrow         int     `json:"years_to_grow"         binding:"required"`
}

// InvestmentGrowthResponse defines the structure for an investment growth calculation response.
type InvestmentGrowthResponse struct {
	FutureValue         float64 `json:"future_value"`
	TotalContributions  float64 `json:"total_contributions"`
	TotalInterestEarned float64 `json:"total_interest_earned"`
}

// RentVsBuyRequest defines the structure for a rent vs. buy calculation request.
type RentVsBuyRequest struct {
	// Buying Inputs
	HomePrice          float64 `json:"home_price"            binding:"required"`
	DownPayment        float64 `json:"down_payment"          binding:"required"`
	InterestRate       float64 `json:"interest_rate"         binding:"required"`
	LoanTermYears      int     `json:"loan_term_years"       binding:"required"`
	PropertyTaxRate    float64 `json:"property_tax_rate"     binding:"required"`
	HomeInsurance      float64 `json:"home_insurance"        binding:"required"` // Annual
	MaintenanceCosts   float64 `json:"maintenance_costs"     binding:"required"` // Annual
	AppreciationRate   float64 `json:"appreciation_rate"     binding:"required"`
	ClosingCosts       float64 `json:"closing_costs"         binding:"required"`
	SellingCostsRate   float64 `json:"selling_costs_rate"    binding:"required"`

	// Renting Inputs
	MonthlyRent        float64 `json:"monthly_rent"          binding:"required"`
	RentersInsurance   float64 `json:"renters_insurance"     binding:"required"` // Annual
	AnnualRentIncrease float64 `json:"annual_rent_increase"  binding:"required"`

	// Shared Inputs
	ComparisonYears     int     `json:"comparison_years"      binding:"required"`
	InvestmentReturnRate float64 `json:"investment_return_rate" binding:"required"`
}

// RentVsBuyResponse defines the structure for a rent vs. buy calculation response.
type RentVsBuyResponse struct {
	TotalCostOfOwning  float64 `json:"total_cost_of_owning"`
	TotalCostOfRenting float64 `json:"total_cost_of_renting"`
	NetBenefitOfOwning float64 `json:"net_benefit_of_owning"`
	Recommendation     string  `json:"recommendation"`
}

// TaxEstimatorRequest defines the structure for a tax estimation request.
type TaxEstimatorRequest struct {
	FilingStatus string  `json:"filing_status" binding:"required"`
	GrossIncome  float64 `json:"gross_income"  binding:"required"`
	Deductions   float64 `json:"deductions"`
}

// TaxEstimatorResponse defines the structure for a tax estimation response.
type TaxEstimatorResponse struct {
	EffectiveTaxRate float64 `json:"effective_tax_rate"`
	TotalTaxAmount   float64 `json:"total_tax_amount"`
}
