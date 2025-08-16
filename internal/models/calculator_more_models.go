package models

// Loan comparison
type LoanOfferModel struct {
	Principal    float64 `json:"principal" binding:"required"`
	InterestRate float64 `json:"interest_rate" binding:"required"`
	TermYears    int     `json:"term_years" binding:"required"`
	Fees         float64 `json:"fees"`
}

type LoanComparisonRequest struct {
	Offers []LoanOfferModel `json:"offers" binding:"required"`
}

type LoanComparisonResponse struct {
	Results []map[string]float64 `json:"results"`
}

// Mortgage affordability
type AffordabilityRequest struct {
	AnnualIncome     float64 `json:"annual_income" binding:"required"`
	DTIRatio         float64 `json:"dti_ratio" binding:"required"`
	DownPayment      float64 `json:"down_payment"`
	InterestRate     float64 `json:"interest_rate" binding:"required"`
	LoanTermYears    int     `json:"loan_term_years" binding:"required"`
	OtherMonthlyDebts float64 `json:"other_monthly_debts"`
}

type AffordabilityResponse struct {
	PurchasePrice float64 `json:"purchase_price"`
	MaxMonthly    float64 `json:"max_monthly"`
}

// Additional simple request/response models
type CreditOptRequest struct {
	Cards []CardInputModel `json:"cards" binding:"required"`
	ExtraMonthly float64 `json:"extra_monthly"`
}

type CreditOptResponse struct {
	Recommendation string `json:"recommendation"`
}

type CollegeSavingsRequest struct {
	CurrentBalance float64 `json:"current_balance"`
	AnnualTuition  float64 `json:"annual_tuition"`
	TuitionInflation float64 `json:"tuition_inflation"`
	AnnualReturn   float64 `json:"annual_return"`
	YearsUntilCollege int `json:"years_until_college"`
}

type CollegeSavingsResponse struct {
	TotalNeed float64 `json:"total_need"`
	RequiredMonthly float64 `json:"required_monthly"`
}

type FeeDragRequest struct {
	Initial float64 `json:"initial"`
	AnnualContribution float64 `json:"annual_contribution"`
	GrossReturn float64 `json:"gross_return"`
	FeePercent float64 `json:"fee_percent"`
	Years int `json:"years"`
}

type FeeDragResponse struct {
	Gross float64 `json:"gross"`
	Net float64 `json:"net"`
	Difference float64 `json:"difference"`
}

type SafeWithdrawalRequest struct {
	Initial float64 `json:"initial"`
	WithdrawalRate float64 `json:"withdrawal_rate"`
	AnnualReturn float64 `json:"annual_return"`
	Years int `json:"years"`
}

type SafeWithdrawalResponse struct {
	Remaining float64 `json:"remaining"`
}

// CD Ladder request/response
type CDLadderRequest struct {
	TotalAmount float64 `json:"total_amount"`
	Rungs int `json:"rungs"`
	BaseRate float64 `json:"base_rate"`
}

type CDLadderResponse struct {
	Maturities []float64 `json:"maturities"`
}

// Payroll & currency
type PayrollRequest struct {
	GrossAnnual float64 `json:"gross_annual"`
	PayPeriods int `json:"pay_periods"`
	FederalRate float64 `json:"federal_rate"`
	StateRate float64 `json:"state_rate"`
	PreTax float64 `json:"pre_tax"`
}

type PayrollResponse struct {
	NetPerPeriod float64 `json:"net_per_period"`
}

type ConvertCurrencyRequest struct {
	Amount float64 `json:"amount"`
	Rate float64 `json:"rate"`
}

type ConvertCurrencyResponse struct {
	Converted float64 `json:"converted"`
}
