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
## API Endpoints
### User
| Method | Endpoint         | Description           |
| ------ | ---------------- | --------------------- |
| POST   | `/user/register` | Register a new user   |
| POST   | `/user/login`    | Login and receive JWT |
### Account
| Method | Endpoint                           | Description                |
| ------ | ---------------------------------- | -------------------------- |
| POST   | `/account/`                        | Create a bank account      |
| POST   | `/account/deposit`                 | Deposit money              |
| POST   | `/account/withdraw`                | Withdraw money             |
| POST   | `/account/transfer`                | Transfer money             |
| GET    | `/account/:accountNo`              | Get account details        |
| PUT    | `/account/:accountNo/status`       | Update account status      |
| GET    | `/account/:accountNo/statistics`   | Get transaction statistics |
| GET    | `/account/:accountNo/summary`      | Get account summary        |
| GET    | `/account/:accountNo/transactions` | Get account transactions   |

All account endpoints require JWT authentication.
## Transaction Filtering
The transaction endpoint supports:
- Pagination
- Transaction type filtering
- Date range filtering
- Sorting
Example:
```text
GET /account/ACC12345678/transactions?page=1&limit=5&type=DEPOSIT&from=2026-08-01&to=2026-08-31&sort=newest
```
Supported transaction types:
```text
DEPOSIT
WITHDRAW
TRANSFER
```
## Database
The application uses MySQL with GORM.
Main Entities: 
```text
User
 │
 └── Account
       │
       └── Transaction
```
Users have accounts, and accounts have transaction records.
Email addresses are protected with a unique database constraint.
## Money Transfer
Money transfers are performed inside a database transaction.
The application locks the relevant account rows while transferring money to prevent concurrent operations from causing inconsistent balances.
```text
Begin Transaction
       │
       ▼
Lock Sender Account
       │
       ▼
Lock Receiver Account
       │
       ▼
Validate Balance & Status
       │
       ▼
Update Balances
       │
       ▼
Create Transaction Records
       │
       ▼
Commit Transaction
```
If an error occurs, the transaction is rolled back.
## Error Handling
This application uses custom error such as:
```text
ErrAccountNotFound
ErrInvalidStatus
ErrInsufficientBalance
ErrInvalidAmount
ErrSameAccountTransfer
ErrEmailAlreadyExists
```
Errors are mapped to appropriate HTTP responses.
Internal server errors are not exposed directly to API clients.
## Swagger Documentation
Swagger is integrated for interactive API documentation.
After starting the application, open:
```text
http://localhost:8080/swagger/index.html
```
Swagger provides:
- Available endpoints
- Request parameters
- Request body examples
- Authentication
- Response documentation
- Interactive API testing
## Getting Started
### Prerequisites

Make sure you have:

- Go installed
- MySQL installed and running
- A MySQL database created
## Clone the repository
```bash
git clone <>
```
