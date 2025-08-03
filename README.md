# Banking Application in Go

A modern banking application built with Go (Golang) that enables management of users, bank accounts, and virtual cards. The system supports core banking operations such as fund transfers, deposits, withdrawals, and secure user authentication.

## 🛠 Technologies and Stack

- **Language**: Go (Golang)
- **APIs**: RESTful (HTTP), gRPC
- **Logging**: [logrus](https://github.com/sirupsen/logrus)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Monitoring**: Prometheus (metrics), OpenTelemetry + Jaeger (tracing)
- **Infrastructure**: Docker, Docker Compose

## 🚀 Features

### 1. User Service
- User registration
- Secure login and authentication (JWT-based)

### 2. Banking Service
- Create new accounts per user
- Deposit and withdraw funds
- Transfer money between accounts (internal and external)

### 3. Card Service
- Generate virtual debit cards
- View card details (masked PAN, expiry, status)
- Block or deactivate cards

## 📌 API Documentation

- **REST API**: OpenAPI 3.0 specification available in `api/rest/`
- **gRPC Services**: Protocol buffer (`.proto`) files in `api/grpc/`

## 🧪 Project Structure

├── api/ # OpenAPI & Protobuf definitions

├── deploy/ # docker-compose

├── configs/ # Configuration files (YAML, env)

├── {service}/

│ ├── cmd/ # Application entry points

│ ├── internal/ # Internal application logic

│ ├── deploy/ # Docker

│ ├── configs/ # Configuration files (YAML, json)

└ └─ pkg/ # Shared utilities (auth, db, middleware, etc.)

## Deploy

With command:

docker-compose -f ./deploy/compose.yml -p "banking_app" up --build app

Default address: localhost:8080
