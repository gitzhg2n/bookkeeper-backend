package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register all route handlers
	RegisterAuthRoutes(r)
	RegisterAccountRoutes(r)
	RegisterBudgetRoutes(r)
	RegisterGoalRoutes(r)
	RegisterInvestmentRoutes(r)
	RegisterTransactionRoutes(r)
	RegisterHouseholdRoutes(r)
	RegisterUserRoutes(r)
	RegisterIncomeSourceRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Bookkeeper backend running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}