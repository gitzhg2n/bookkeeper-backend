# Bookkeeper Deployment Instructions

## Prerequisites
- Docker and Docker Compose installed
- (Optional) Custom domain and SSL for production

## 1. Environment Setup
- Copy `.env.example` to `.env` in the backend directory and set secrets, database info, and email/push providers as needed.
- (Optional) Adjust `REACT_APP_API_BASE` in the frontend `.env` if deploying backend separately or using a custom domain.

## 2. Build and Run All Services
From the `bookkeeper-backend` directory:

```sh
docker-compose up --build
```

- Backend: http://localhost:3000
- Frontend: http://localhost:8080
- Database: Internal (Postgres)

## 3. Production Tips
- Set strong secrets in `.env` (never use defaults in production)
- Use a secure Postgres instance (managed or self-hosted)
- Set up HTTPS (use a reverse proxy like Nginx or Caddy)
- Configure email/push providers for notifications
- Use a custom domain for frontend and backend

## 4. Updating
To update the app, pull the latest code and re-run:
```sh
git pull
# (in both backend and frontend folders)
docker-compose up --build
```

## 5. Stopping Services
```sh
docker-compose down
```

---
For advanced deployment (Kubernetes, cloud, etc.), see the documentation or contact the maintainers.
