# Banking API
A RESTful Banking API built with Go, Gin, GORM, and MySQL.  
The application provides user authentication, bank account management, transactions, money transfers, transaction history, statistics, and account summaries.
## Features
- User registration and login
- JWT-based authentication
- Create and manage bank accounts
- Deposit money
- Withdraw money
- Transfer money between accounts
- Account status management
- Account details
- Transaction history
- Transaction filtering and pagination
- Transaction statistics
- Account financial summary
- Swagger API documentation
- Password hashing with bcrypt
- Database transactions for money transfers
- Custom application errors
- Layered architecture
## Tech Stack
- **Go** — Backend programming language
- **Gin** — HTTP web framework for building REST APIs
- **GORM** — ORM for interacting with MySQL
- **MySQL** — Relational database
- **JWT** — Token-based authentication
- **bcrypt** — Secure password hashing
- **Swagger** — API documentation
- **Docker** — Application containerization
## Architecture
This application follows a layered architecture:
```text
Client
  │
  ▼
Controller
  │
  ▼
Service
  │
  ▼
Repository
  │
  ▼
MySQL Database
```
### Layers

### Controller

Handles HTTP requests and responses
Validates request input
Maps application errors to HTTP status codes

### Service

Contains business logic
Handles banking operations such as deposits, withdrawals, and transfers

### Repository

Handles database operations
Uses GORM to communicate with MySQL
## Project Structure
```text
bankingApp/
│
├── config/
│   ├── config.go
│   └── database.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── errors/
│   └── errors.go
│
├── internal/
│   ├── account/
│   │   ├── account_controller.go
│   │   ├── account_service.go
│   │   ├── account_repository.go
│   │   ├── dto.go
│   │   └── module.go
│   │
│   └── user/
│       ├── user_controller.go
│       ├── user_service.go
│       ├── user_repository.go
│       ├── dto.go
│       └── module.go
│
├── middleware/
│   └── auth.go
│
├── models/
│   ├── user.go
│   ├── account.go
│   └── transaction.go
│
├── utils/
│   └── response/
│
├── main.go
├── go.mod
└── go.sum
```
## Authentication
The API uses JWT-based authentication.
### Authentication flow
```text
Register
   │
   ▼
Password hashed with bcrypt
   │
   ▼
User stored in MySQL
   │
   ▼
Login
   │
   ▼
JWT generated
   │
   ▼
JWT sent with protected requests
```
Protected endpoints require:
```text
Authorization: Bearer <token>
```
