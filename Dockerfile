# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o bookkeeper-backend ./main.go

# ---- Runtime Stage ----
FROM alpine:3.20

WORKDIR /app

# Copy binary and static files from builder
COPY --from=builder /app/bookkeeper-backend .
COPY --from=builder /app/.env .env

EXPOSE 3000

CMD ["./bookkeeper-backend"]