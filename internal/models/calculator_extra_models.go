package models

import "time"

// Amortization
type AmortizationRequest struct {
    LoanAmount     float64   `json:"loan_amount" binding:"required"`
    InterestRate   float64   `json:"interest_rate" binding:"required"`
    LoanTermYears  int       `json:"loan_term_years" binding:"required"`
    StartDate      time.Time `json:"start_date"`
    ExtraMonthly   float64   `json:"extra_monthly"`
}

type AmortizationResponse struct {
    Schedule      []map[string]interface{} `json:"schedule"`
    TotalInterest float64                  `json:"total_interest"`
}

// Refinance
type RefinanceRequest struct {
    CurrentBalance       float64 `json:"current_balance" binding:"required"`
    CurrentRate          float64 `json:"current_rate" binding:"required"`
    CurrentRemainingYears int    `json:"current_remaining_years" binding:"required"`
    NewRate              float64 `json:"new_rate" binding:"required"`
    NewTermYears         int     `json:"new_term_years" binding:"required"`
    ClosingCosts         float64 `json:"closing_costs"`
    YearsToEvaluate      int     `json:"years_to_evaluate"`
}

type RefinanceResponse struct {
    NewMonthly       float64 `json:"new_monthly"`
    MonthlySavings   float64 `json:"monthly_savings"`
    BreakevenMonths  int     `json:"breakeven_months"`
    TotalSavings     float64 `json:"total_savings"`
}

// APR/APY
type APRAPYRequest struct {
    Rate               float64 `json:"rate" binding:"required"`
    CompoundsPerYear   int     `json:"compounds_per_year" binding:"required"`
}

type APRAPYResponse struct {
    Converted float64 `json:"converted"`
}

// Retirement
type RetirementRequest struct {
    CurrentBalance        float64 `json:"current_balance"`
    AnnualSalary          float64 `json:"annual_salary" binding:"required"`
    ContributionRate      float64 `json:"contribution_rate" binding:"required"`
    EmployerMatchRate     float64 `json:"employer_match_rate"`
    EmployerMatchCapPercent float64 `json:"employer_match_cap_percent"`
    AnnualReturn          float64 `json:"annual_return" binding:"required"`
    Inflation             float64 `json:"inflation"`
    Years                 int     `json:"years" binding:"required"`
}

type RetirementResponse struct {
    NominalFutureValue float64   `json:"nominal_future_value"`
    RealFutureValue    float64   `json:"real_future_value"`
    YearlyProjections  []float64 `json:"yearly_projections"`
}

// Savings goal
type SavingsGoalRequest struct {
    GoalAmount        float64 `json:"goal_amount" binding:"required"`
    CurrentSavings    float64 `json:"current_savings"`
    MonthlyContribution float64 `json:"monthly_contribution"`
    AnnualReturn      float64 `json:"annual_return"`
}

type SavingsGoalResponse struct {
    MonthsToGoal         int     `json:"months_to_goal"`
    RequiredMonthly      float64 `json:"required_monthly"`
}

// Credit card payoff
type CardInputModel struct {
    Balance    float64 `json:"balance" binding:"required"`
    APR        float64 `json:"apr" binding:"required"`
    MinPayment float64 `json:"min_payment" binding:"required"`
}

type CreditPayoffRequest struct {
    Cards       []CardInputModel `json:"cards" binding:"required"`
    ExtraMonthly float64         `json:"extra_monthly"`
    Method      string           `json:"method"` // "snowball" or "avalanche"
}

type CreditPayoffResponse struct {
    MonthsToPayOff int     `json:"months_to_pay_off"`
    TotalInterest  float64 `json:"total_interest"`
    Order          []int   `json:"order"`
}

// Take-home pay
type TakeHomeRequest struct {
    GrossAnnual        float64 `json:"gross_annual" binding:"required"`
    PayPeriodsPerYear  int     `json:"pay_periods_per_year" binding:"required"`
    FederalTaxRate     float64 `json:"federal_tax_rate"`
    StateTaxRate       float64 `json:"state_tax_rate"`
    PreTaxDeductions   float64 `json:"pre_tax_deductions"`
}

type TakeHomeResponse struct {
    NetPerPeriod float64 `json:"net_per_period"`
}

// Inflation adjust
type InflationRequest struct {
    Nominal        float64 `json:"nominal" binding:"required"`
    InflationRate  float64 `json:"inflation_rate" binding:"required"`
    Years          int     `json:"years" binding:"required"`
}

type InflationResponse struct {
    Real float64 `json:"real"`
}

// Net worth
type NetWorthRequest struct {
    Assets     map[string]float64 `json:"assets"`
    Liabilities map[string]float64 `json:"liabilities"`
}

type NetWorthResponse struct {
    NetWorth float64            `json:"net_worth"`
    Allocation map[string]float64 `json:"allocation"`
}
