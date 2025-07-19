# Bookkeeper Backend (Go Version)

This repository contains the backend code for **Bookkeeper**, a privacy-first, open-source personal finance management application written in Go.

## 🚀 Quick Start

### Prerequisites
- Go 1.22 or later
- Docker (optional)

### Local Development

1. **Clone and setup**:
```bash
git clone <repository-url>
cd bookkeeper-backend
cp .env.example .env
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Run the application**:
```bash
go run main.go
```

The server will start on port 3000 (or the port specified in your `.env` file).

### Using Docker

1. **Build the Docker image**:
```bash
docker build -t bookkeeper-backend .
```

2. **Run with Docker Compose**:
```bash
docker-compose up
```

## 📋 API Endpoints

### Health Checks
- `GET /health` - Basic health check
- `GET /ready` - Readiness check (includes database connectivity)

### Financial Calculators
- `POST /calculators/mortgage` - Mortgage payment calculator
- `POST /calculators/rent-vs-buy` - Rent vs buy comparison
- `POST /calculators/investment-growth` - Investment growth projector
- `POST /calculators/debt-payoff` - Debt payoff calculator
- `POST /calculators/tax-estimator` - Tax estimation

### Core API
- `/auth` - Authentication endpoints
- `/accounts` - Account management
- `/budgets` - Budget planning
- `/goals` - Financial goals
- `/investments` - Investment tracking
- `/transactions` - Transaction management
- `/households` - Household management
- `/users` - User management
- `/incomeSources` - Income source tracking

## 🏛️ Architecture Overview

### Backend Architecture
