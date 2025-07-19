package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookkeeper-backend/models"
	"bookkeeper-backend/routes"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	models.InitDB()

	r := mux.NewRouter()

	// Add health check endpoint
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/ready", readinessCheck).Methods("GET")

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

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Bookkeeper backend running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"bookkeeper-backend"}`))
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	sqlDB, err := models.DB.DB()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"not ready","error":"database connection failed"}`))
		return
	}
	
	if err := sqlDB.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"not ready","error":"database ping failed"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","service":"bookkeeper-backend"}`))
}