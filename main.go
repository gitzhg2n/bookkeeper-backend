package main

import (
	"log"
	"net/http"
	"os"

	"bookkeeper-backend/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register all route handlers
	routes.RegisterAuthRoutes(r)
	routes.RegisterAccountRoutes(r)
	routes.RegisterBudgetRoutes(r)
	routes.RegisterGoalRoutes(r)
	routes.RegisterInvestmentRoutes(r)
	routes.RegisterTransactionRoutes(r)
	routes.RegisterHouseholdRoutes(r)
	routes.RegisterHouseholdMemberRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisterIncomeSourceRoutes(r)
	routes.RegisterCalculatorRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Bookkeeper backend running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}