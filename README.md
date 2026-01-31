# Personal Monthly Expense Tracker - Backend

A clean, stateless Go backend for tracking personal monthly expenses with JWT authentication and MongoDB.

## Features

- ✅ User authentication (sign up, login) with JWT tokens
- ✅ Monthly budget management
- ✅ Base income tracking
- ✅ Expense management (add, update, delete)
- ✅ Automatic budget creation for new months
- ✅ Remaining balance calculation
- ✅ User data isolation (users only see their own data)

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: MongoDB
- **Authentication**: JWT (Bearer tokens)
- **Password Hashing**: bcrypt

## Project Structure

```
finance-app/
├── cmd/server/           # Application entry point
├── internal/
│   ├── auth/             # JWT and middleware
│   ├── config/           # Configuration management
│   ├── db/               # MongoDB connection
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   └── utils/            # Helper functions
├── go.mod
└── .env                  # Environment variables
```

## Prerequisites

- Go 1.21 or higher
- MongoDB 4.0 or higher
- Make (optional)

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd finance-app
```

### 2. Install dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure environment variables

Copy `.env.example` to `.env` and update values:

```bash
cp .env.example .env
```

Edit `.env`:

```
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=finance_app
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
PORT=3000
```

### 4. Start MongoDB

```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or using local MongoDB installation
mongod
```

### 5. Run the server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:3000`

## API Endpoints

### Authentication

#### Sign Up

```
POST /auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response: 201 Created
{
  "token": "eyJhbGc..."
}
```

#### Login

```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}

Response: 200 OK
{
  "token": "eyJhbGc..."
}
```

### Budget Management

All budget endpoints require authentication: `Authorization: Bearer <token>`

#### Get Current Month's Budget

```
GET /budget/current

Response: 200 OK
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Rent",
      "amount": 1200,
      "createdAt": "2026-01-31T10:30:00Z"
    }
  ],
  "remaining": 3800
}
```

#### Get Specific Month's Budget

```
GET /budget?year=2026&month=1

Response: 200 OK
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [...],
  "remaining": 3800
}
```

#### Set Base Income

```
POST /budget/base-income
Content-Type: application/json

{
  "amount": 5000
}

Response: 200 OK
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [],
  "remaining": 5000
}
```

### Expense Management

All expense endpoints require authentication: `Authorization: Bearer <token>`

#### Add Expense

```
POST /expenses
Content-Type: application/json

{
  "title": "Rent",
  "amount": 1200
}

Response: 201 Created
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Rent",
      "amount": 1200,
      "createdAt": "2026-01-31T10:30:00Z"
    }
  ],
  "remaining": 3800
}
```

#### Update Expense

```
PUT /expenses/:expenseId
Content-Type: application/json

{
  "title": "Updated Rent",
  "amount": 1300
}

Response: 200 OK
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [...],
  "remaining": 3700
}
```

#### Delete Expense

```
DELETE /expenses/:expenseId

Response: 200 OK
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [],
  "remaining": 5000
}
```

## Authentication

All protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

The token is valid for the duration specified in `JWT_EXPIRY_HOURS` (default: 24 hours).

## Domain Rules

### Users

- Email must be unique
- Password is hashed with bcrypt
- Users can only access their own data

### Monthly Budgets

- One budget per user per month (identified by userId + year + month)
- Automatically created when accessed if doesn't exist
- Base income is optional (can be null)

### Expenses

- Belong to a specific user and month
- Must have title and positive amount
- Can be added, updated, and deleted
- Expense IDs are unique per budget

### Remaining Balance

- Calculated as: `baseIncome - sum(expenses.amount)`
- Returns `null` if base income is not set
- Not stored in database (derived value)

## Error Handling

All endpoints return JSON error responses:

```json
{
  "error": "descriptive error message"
}
```

Common HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o finance-app cmd/server/main.go
```

## Deployment

The application is designed to be deployment-ready for:

- Docker (containerization)
- Fly.io
- Render
- Vercel (with refactoring)

No special deployment code is included at this stage. Follow standard Go deployment practices for your chosen platform.

## License

MIT License

## Support

For issues or questions, please open an issue in the repository.
