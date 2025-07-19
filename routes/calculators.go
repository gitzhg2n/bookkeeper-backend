package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"github.com/gorilla/mux"
)

// --- Calculator Payloads ---

type RentBuyPayload struct {
	RentAmount      float64 `json:"rentAmount"`
	HomePrice       float64 `json:"homePrice"`
	DownPayment     float64 `json:"downPayment"`
	InterestRate    float64 `json:"interestRate"`
	LoanTermYears   int     `json:"loanTermYears"`
	PropertyTaxRate float64 `json:"propertyTaxRate"`
	InsuranceRate   float64 `json:"insuranceRate"`
	AnnualIncrease  float64 `json:"annualIncrease"`
}

type MortgagePayload struct {
	HomePrice     float64 `json:"homePrice"`
	DownPayment   float64 `json:"downPayment"`
	InterestRate  float64 `json:"interestRate"`
	LoanTermYears int     `json:"loanTermYears"`
}

type InvestmentGrowthPayload struct {
	InitialAmount float64 `json:"initialAmount"`
	AnnualRate    float64 `json:"annualRate"`
	Years         int     `json:"years"`
}

type DebtPayoffPayload struct {
	Balance          float64 `json:"balance"`
	AnnualRate       float64 `json:"annualRate"`
	MonthlyPayment   float64 `json:"monthlyPayment"`
	ExtraPerMonth    float64 `json:"extraPerMonth"`
}

type TaxEstimatorPayload struct {
	Income float64 `json:"income"`
	State  string  `json:"state"`
}

// --- Route Registrar ---
func RegisterCalculatorRoutes(r *mux.Router) {
	sub := r.PathPrefix("/calculators").Subrouter()
	sub.HandleFunc("/rent-vs-buy", rentVsBuyCalculator).Methods("POST")
	sub.HandleFunc("/mortgage", mortgageCalculator).Methods("POST")
	sub.HandleFunc("/investment-growth", investmentGrowthCalculator).Methods("POST")
	sub.HandleFunc("/debt-payoff", debtPayoffCalculator).Methods("POST")
	sub.HandleFunc("/tax-estimator", taxEstimatorCalculator).Methods("POST")
	// Add more endpoints as needed
}

// --- Calculator Implementations ---

func rentVsBuyCalculator(w http.ResponseWriter, r *http.Request) {
	var payload RentBuyPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Simple calculation: Total cost of renting vs buying over loan term
	rentTotal := payload.RentAmount * float64(payload.LoanTermYears) * 12
	loanPrincipal := payload.HomePrice - payload.DownPayment
	monthlyRate := payload.InterestRate / 100 / 12
	nPayments := payload.LoanTermYears * 12
	// Mortgage payment calculation
	monthlyMortgage := 0.0
	if monthlyRate > 0 {
		monthlyMortgage = loanPrincipal * (monthlyRate * math.Pow(1+monthlyRate, float64(nPayments))) / (math.Pow(1+monthlyRate, float64(nPayments)) - 1)
	} else {
		monthlyMortgage = loanPrincipal / float64(nPayments)
	}
	buyTotal := monthlyMortgage * float64(nPayments)
	// Add property tax & insurance
	buyTotal += payload.HomePrice * (payload.PropertyTaxRate/100 + payload.InsuranceRate/100) * float64(payload.LoanTermYears)

	result := map[string]interface{}{
		"rentTotal": rentTotal,
		"buyTotal":  buyTotal,
		"monthlyMortgage": monthlyMortgage,
		"inputs":    payload,
	}
	json.NewEncoder(w).Encode(result)
}

func mortgageCalculator(w http.ResponseWriter, r *http.Request) {
	var payload MortgagePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	loanPrincipal := payload.HomePrice - payload.DownPayment
	monthlyRate := payload.InterestRate / 100 / 12
	nPayments := payload.LoanTermYears * 12
	monthlyPayment := 0.0
	if monthlyRate > 0 {
		monthlyPayment = loanPrincipal * (monthlyRate * math.Pow(1+monthlyRate, float64(nPayments))) / (math.Pow(1+monthlyRate, float64(nPayments)) - 1)
	} else {
		monthlyPayment = loanPrincipal / float64(nPayments)
	}
	// Amortization schedule omitted for brevity
	result := map[string]interface{}{
		"monthlyPayment": monthlyPayment,
		"totalPaid":      monthlyPayment * float64(nPayments),
		"inputs":         payload,
	}
	json.NewEncoder(w).Encode(result)
}

func investmentGrowthCalculator(w http.ResponseWriter, r *http.Request) {
	var payload InvestmentGrowthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Compound interest formula
	fv := payload.InitialAmount * math.Pow(1+payload.AnnualRate/100, float64(payload.Years))
	result := map[string]interface{}{
		"futureValue": fv,
		"inputs":      payload,
	}
	json.NewEncoder(w).Encode(result)
}

func debtPayoffCalculator(w http.ResponseWriter, r *http.Request) {
	var payload DebtPayoffPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Simple debt payoff: amortized payoff date (ignoring interest changes, snowball, etc.)
	balance := payload.Balance
	rate := payload.AnnualRate / 100 / 12
	payment := payload.MonthlyPayment + payload.ExtraPerMonth
	months := 0
	totalInterest := 0.0
	for balance > 0 && months < 600 {
		interest := balance * rate
		balance += interest
		balance -= payment
		totalInterest += interest
		months++
		if payment <= interest {
			break // Will never pay off
		}
	}
	result := map[string]interface{}{
		"monthsToPayoff": months,
		"totalInterest":  totalInterest,
		"inputs":         payload,
	}
	json.NewEncoder(w).Encode(result)
}

func taxEstimatorCalculator(w http.ResponseWriter, r *http.Request) {
	var payload TaxEstimatorPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Simple US federal tax brackets (2025), single filer (ignore state tax for now)
	income := payload.Income
	tax := 0.0
	brackets := []struct {
		limit float64
		rate  float64
	}{
		{11000, 0.10},
		{44725, 0.12},
		{95375, 0.22},
		{182100, 0.24},
		{231250, 0.32},
		{578125, 0.35},
	}
	lastLimit := 0.0
	for _, b := range brackets {
		if income > b.limit {
			tax += (b.limit - lastLimit) * b.rate
			lastLimit = b.limit
		} else {
			tax += (income - lastLimit) * b.rate
			lastLimit = income
			break
		}
	}
	if income > 578125 {
		tax += (income - 578125) * 0.37
	}
	result := map[string]interface{}{
		"estimatedTax": tax,
		"inputs":       payload,
	}
	json.NewEncoder(w).Encode(result)
}