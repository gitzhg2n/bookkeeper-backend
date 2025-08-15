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


## 📚 API Details

### Authentication
- `POST /auth/login` — Log in and receive a JWT token
- `POST /auth/signup` — Create a new user
- `POST /auth/logout` — Log out
- `POST /auth/password-reset` — Request password reset

### Accounts
- `GET /accounts` — List user accounts
- `POST /accounts` — Create account
- `GET /accounts/{id}` — Get account details
- `PUT /accounts/{id}` — Update account
- `DELETE /accounts/{id}` — Delete account

### Transactions
- `GET /transactions?account_id=...` — List transactions for an account
- `POST /transactions/{account_id}` — Create transaction
- `PUT /transactions/{id}` — Update transaction
- `DELETE /transactions/{id}` — Delete transaction

### Budgets
- `GET /budgets?household_id=...&month=YYYY-MM` — List budgets
- `POST /budgets` — Create or update budget

### Categories
- `GET /categories?household_id=...` — List categories
- `POST /categories` — Create category
- `PUT /categories/{id}` — Update category
- `DELETE /categories/{id}` — Delete category

### Goals
- `GET /goals` — List financial goals
- `POST /goals` — Create goal
- `PUT /goals/{id}` — Update goal
- `DELETE /goals/{id}` — Delete goal

### Notifications & Alerts
- `GET /notifications` — List notifications
- `POST /notifications/{id}/read` — Mark notification as read
- `POST /notifications/read_all` — Mark all as read
- `GET /investment_alerts` — List investment alerts
- `POST /investment_alerts` — Create alert
- `PUT /investment_alerts/{id}` — Update alert
- `DELETE /investment_alerts/{id}` — Delete alert

### User Settings
- `GET /user_settings` — Get user settings
- `PUT /user_settings` — Update user settings (notification preferences, etc)

### Households
- `GET /households` — List households
- `POST /households` — Create household
- `PUT /households/{id}` — Update household
- `DELETE /households/{id}` — Delete household

### Example: Create Transaction
```http
POST /transactions/{account_id}
{
	"amount_cents": 12345,
	"currency": "USD",
	"category_id": 1,
	"memo": "Groceries",
	"occurred_at": "2025-08-15T12:00:00Z"
}
```

## 🧩 Main Models

- **User**: Authentication, profile, and plan info
- **Account**: Financial account (bank, credit, etc)
- **Transaction**: Linked to account, category, and user
- **Budget**: Monthly planned spending per category
- **Category**: User-defined or default spending categories
- **Goal**: Financial goals (amount, due date)
- **Notification**: In-app and email/push notifications
- **InvestmentAlert**: Customizable investment/price alerts
- **AlertHistory**: Tracks alert triggers and cooldowns
- **UserSettings**: Notification preferences, premium features

## 🚨 Advanced Features

- **Notifications**: In-app, email, and push. User-configurable preferences.
- **Investment Alerts**: Compound, time-based, and custom rule logic. Cooldowns and alert history for premium/self-hosted users.
- **Premium/Self-Hosted**: Advanced alert rules, unlimited notifications, custom thresholds, and privacy-first design.
- **Security**: End-to-end encryption for sensitive data, strong password hashing, and key management.

## 🚀 Deployment

### Docker Compose (Recommended)
1. Copy `.env.example` to `.env` and set secrets/DB info.
2. Run:
	 ```bash
	 docker-compose up --build
	 ```
3. The backend will be available at `http://localhost:3000` (or your configured port).

### Production Tips
- Use a secure database (Postgres recommended)
- Set strong secrets in `.env`
- Use HTTPS in production
- Set up email/push providers for notifications

### Contributing
- See `CONTRIBUTING.md` for guidelines
- Run tests before submitting PRs

---
For full API details, see the code in `/routes` and `/internal/models`. For advanced alert logic, see `/internal/jobs/investment_alerts.go`.
