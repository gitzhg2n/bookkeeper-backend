# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o bookkeeper-backend ./cmd/server

# ---- Runtime Stage ----
FROM alpine:3.20

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bookkeeper-backend .

EXPOSE 3000

CMD ["./bookkeeper-backend"]