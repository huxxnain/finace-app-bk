# Architecture Documentation

## Overview

The Finance App is a clean, stateless backend built with Go and Fiber framework. It implements a layered architecture with clear separation of concerns.

## Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Handlers                     │
│         (auth, budget, expense handlers)            │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                   Services                          │
│      (business logic, validation, calculations)     │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  Data Layer                         │
│        (MongoDB operations, persistence)            │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│              MongoDB Database                       │
│         (collections, indexes, documents)           │
└─────────────────────────────────────────────────────┘
```

## Directory Structure

```
finance-app/
│
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── auth/                       # Authentication & JWT
│   │   ├── jwt.go                  # Token generation/validation
│   │   └── middleware.go           # Authentication middleware
│   │
│   ├── config/                     # Configuration management
│   │   └── config.go               # Load env variables
│   │
│   ├── db/                         # Database connection
│   │   └── mongo.go                # MongoDB client
│   │
│   ├── handlers/                   # HTTP request handlers
│   │   ├── auth_handler.go         # Sign up, login
│   │   ├── budget_handler.go       # Budget CRUD
│   │   └── expense_handler.go      # Expense CRUD
│   │
│   ├── models/                     # Data structures
│   │   └── models.go               # User, Budget, Expense models
│   │
│   ├── services/                   # Business logic
│   │   ├── user_service.go         # User operations
│   │   └── budget_service.go       # Budget & expense operations
│   │
│   └── utils/                      # Helper functions
│       ├── password.go             # Password hashing
│       └── time.go                 # Date/time utilities
│
├── Dockerfile                      # Docker image definition
├── docker-compose.yml              # Multi-container setup
├── Makefile                        # Build commands
├── go.mod                          # Go module definition
├── .env.example                    # Environment template
├── README.md                       # Project documentation
├── API.md                          # API reference
├── QUICKSTART.md                   # Quick start guide
└── ARCHITECTURE.md                 # This file
```

## Component Details

### 1. Handlers Layer (`internal/handlers/`)

Responsible for:

- HTTP request parsing
- Input validation
- Calling service methods
- Formatting responses
- HTTP status codes

**Files:**

- `auth_handler.go` - Sign up and login endpoints
- `budget_handler.go` - Budget retrieval and base income
- `expense_handler.go` - Expense CRUD operations

**Example:** Handling a request

```
HTTP Request → Handler → Parse JSON → Validate → Service Call → Format Response → HTTP Response
```

### 2. Services Layer (`internal/services/`)

Responsible for:

- Business logic implementation
- Data validation
- Calculations (remaining balance)
- Database operations coordination
- Error handling

**Files:**

- `user_service.go` - User registration and authentication
- `budget_service.go` - Budget and expense operations

**Isolation:** Users only see their own data through userID checks.

### 3. Models Layer (`internal/models/`)

Defines data structures:

- `User` - User account information
- `MonthlyBudget` - Monthly budget with expenses
- `Expense` - Individual expense entry
- `*Request` - HTTP request DTOs
- `*Response` - HTTP response DTOs
- `JWTClaims` - JWT payload

All models include BSON tags for MongoDB serialization.

### 4. Database Layer (`internal/db/`)

MongoDB operations:

- Connection management
- Client reuse
- Database/collection access

**Indexes:**

- `users.email` - Unique constraint
- `monthly_budgets.userId + year + month` - Unique constraint

### 5. Authentication (`internal/auth/`)

JWT implementation:

- Token generation with HS256
- Token validation and parsing
- Middleware for protected routes
- Claims extraction

**Security:**

- Passwords hashed with bcrypt (cost 10)
- JWT secret stored in environment
- Token expiration (configurable)

### 6. Configuration (`internal/config/`)

Environment management:

- Load from `.env` file
- Provide sensible defaults
- Type-safe access

### 7. Utilities (`internal/utils/`)

Helper functions:

- Password hashing and verification
- Current month/year calculation
- Time helpers

---

## Data Flow Examples

### Sign Up Flow

```
1. Client: POST /auth/signup {email, password}
2. Handler: Parse request, validate input
3. Service: Check email uniqueness, hash password
4. Database: Insert user document
5. Service: Generate JWT token
6. Handler: Return token
7. Client: Receives token, stores locally
```

### Add Expense Flow

```
1. Client: POST /expenses {title, amount} + Bearer token
2. Middleware: Validate token, extract userID
3. Handler: Parse request, validate input
4. Service: Get or create budget for current month
5. Service: Add expense to budget
6. Database: Update MongoDB document
7. Service: Calculate remaining balance
8. Handler: Return updated budget
9. Client: Receives budget with new expense
```

### Authorization Flow

```
1. Client: Includes Authorization header with token
2. Middleware: Extract token from "Bearer <token>"
3. Auth: Verify JWT signature and expiration
4. Auth: Extract userID from claims
5. Handler: Gets userID from context
6. Service: Uses userID to filter data
7. Database: Query only user's documents
```

---

## Database Schema

### Users Collection

```javascript
{
  _id: ObjectId,
  email: String (unique),
  password: String (hashed),
  createdAt: ISODate
}
```

### Monthly Budgets Collection

```javascript
{
  _id: ObjectId,
  userId: ObjectId (FK to users),
  year: Number,
  month: Number,
  baseIncome: Number | null,
  expenses: [
    {
      _id: ObjectId,
      title: String,
      amount: Number,
      createdAt: ISODate
    }
  ],
  createdAt: ISODate,
  updatedAt: ISODate
}
```

**Indexes:**

```javascript
// Users
db.users.createIndex({ email: 1 }, { unique: true });

// Monthly Budgets
db.monthly_budgets.createIndex(
  {
    userId: 1,
    year: 1,
    month: 1,
  },
  { unique: true },
);
```

---

## API Request/Response Flow

### Typical Request Lifecycle

```
Request (HTTP)
    ↓
Fiber Router
    ↓
Middleware (logging, CORS if added)
    ↓
Auth Middleware (extract & validate token)
    ↓
Handler (parse, validate)
    ↓
Service (business logic)
    ↓
Database (MongoDB)
    ↓
Service (format response)
    ↓
Handler (HTTP status)
    ↓
Response (JSON)
```

### Error Handling

Errors flow back up the stack:

```
Database Error
    ↓
Service catches & wraps error
    ↓
Handler catches error
    ↓
Handler returns error response with status code
    ↓
Client receives error JSON
```

---

## Security Considerations

### 1. Authentication

- JWT with HS256 signing
- Bearer token in Authorization header
- Token expiration (default 24 hours)

### 2. Password Security

- Bcrypt hashing (cost 10)
- Never stored in plaintext
- Never returned in responses

### 3. Authorization

- User can only access own data
- userID from token compared with document owner
- No admin role needed (single user system)

### 4. Data Isolation

- MongoDB queries filtered by userID
- Budget documents linked to userID
- Expenses nested within budget documents

### 5. Input Validation

- Email format validation
- Amount range validation (non-negative)
- Required field validation

### 6. Environment Secrets

- JWT_SECRET stored in `.env`
- MongoDB credentials in connection URI
- Different values for dev/prod

---

## Scalability & Performance

### Current Design

- Stateless handlers (scale horizontally)
- Connection pooling (MongoDB client reused)
- Single-user per token (no session management)
- Embedded expenses (no additional queries)

### MongoDB Efficiency

- Embedded expenses in budget document
  - Single query per budget operation
  - No joins needed
  - Good for small expense lists

- Unique index on (userId, year, month)
  - Fast lookups
  - Prevents duplicates
  - Efficient updates

### Future Optimizations

- Add caching layer (Redis)
- Implement pagination for expense lists
- Add rate limiting
- Connection pooling configuration
- Query result caching

---

## Testing Strategy

### Unit Tests (to implement)

- Service business logic
- Utility functions
- Password hashing

### Integration Tests (to implement)

- Full request/response cycles
- Database operations
- Auth middleware

### Load Testing (optional)

- User signup/login
- Budget operations
- Expense CRUD

---

## Deployment Architecture

### Docker Deployment

```
┌─────────────────────────────────────────┐
│        Docker Container Registry         │
└──────────────────────┬──────────────────┘
                       │
       ┌───────────────┴───────────────┐
       ▼                               ▼
┌─────────────────┐          ┌─────────────────┐
│  Finance App    │          │    MongoDB      │
│  Container      │◄────────►│   Container     │
│  Port 3000      │          │   Port 27017    │
└─────────────────┘          └─────────────────┘
```

### Docker Compose Setup

- Multi-container orchestration
- Network isolation
- Volume persistence
- Health checks

### Environment Variables

- Separate configs for dev/test/prod
- Secrets stored outside code
- No hardcoded credentials

---

## Dependency Management

### Go Modules

- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/golang-jwt/jwt/v5` - JWT handling
- `go.mongodb.org/mongo-driver` - Database driver
- `golang.org/x/crypto` - Password hashing
- `github.com/joho/godotenv` - Environment loading

---

## Best Practices Implemented

✅ **Separation of Concerns**

- Handlers → HTTP only
- Services → Business logic
- Models → Data structures
- DB → Persistence

✅ **Error Handling**

- Explicit error checks
- Meaningful error messages
- Appropriate HTTP status codes

✅ **Security**

- Password hashing
- JWT authentication
- Data isolation
- Input validation

✅ **Code Organization**

- Clear directory structure
- Logical grouping
- Reusable utilities
- Dependency injection via constructors

✅ **Statelessness**

- No session state
- Token-based auth
- Horizontal scalability

---

## Future Enhancements

### Phase 2

- [ ] Unit and integration tests
- [ ] API rate limiting
- [ ] Request logging/monitoring
- [ ] Pagination for large datasets
- [ ] Category management for expenses
- [ ] Recurring expenses

### Phase 3

- [ ] Notifications (email/push)
- [ ] Multi-currency support
- [ ] Budget alerts
- [ ] Export to CSV/PDF
- [ ] Advanced analytics

### Phase 4

- [ ] Multi-user households
- [ ] Shared budgets
- [ ] Admin dashboard
- [ ] Mobile app integration

---

## Related Documentation

- [README.md](README.md) - Project overview
- [API.md](API.md) - API reference
- [QUICKSTART.md](QUICKSTART.md) - Getting started
