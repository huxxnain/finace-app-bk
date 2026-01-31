# API Documentation

## Base URL

```
http://localhost:3000
```

## Response Format

All responses are JSON. Every response contains either `data` and `error` fields (errors) or the actual response body.

## HTTP Status Codes

- `200 OK` - Successful GET/PUT/DELETE
- `201 Created` - Successful POST (resource created)
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

---

## Authentication Endpoints

### POST /auth/signup

Sign up a new user.

**Request**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (201 Created)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**

- `400` - Missing email or password
- `409` - User with this email already exists

**Rules**

- Email must be unique
- Password is hashed with bcrypt
- Token is valid for 24 hours (configurable)

---

### POST /auth/login

Log in an existing user.

**Request**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (200 OK)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors**

- `400` - Missing email or password
- `401` - Invalid email or password

---

## Budget Endpoints

### GET /budget/current

Retrieve the current month's budget. Creates it if it doesn't exist.

**Headers**

```
Authorization: Bearer <token>
```

**Response** (200 OK)

```json
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
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Groceries",
      "amount": 300,
      "createdAt": "2026-01-31T11:00:00Z"
    }
  ],
  "remaining": 3500
}
```

**Errors**

- `401` - Missing or invalid token

**Rules**

- Creates budget automatically if it doesn't exist
- `remaining` is null if `baseIncome` is null
- Budget is unique per user per month

---

### GET /budget

Retrieve a specific month's budget. Creates it if it doesn't exist.

**Headers**

```
Authorization: Bearer <token>
```

**Query Parameters**

- `year` (required): YYYY format (e.g., 2026)
- `month` (required): 1-12 (January to December)

**Example**

```
GET /budget?year=2026&month=1
```

**Response** (200 OK)

```json
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [],
  "remaining": 5000
}
```

**Errors**

- `400` - Missing year or month parameters, invalid format, or month not 1-12
- `401` - Missing or invalid token

**Rules**

- Creates budget automatically if it doesn't exist
- Can retrieve past or future months
- `remaining` is null if `baseIncome` is null

---

### POST /budget/base-income

Set or update the base income for the current month.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**

```json
{
  "amount": 5000
}
```

**Response** (200 OK)

```json
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [],
  "remaining": 5000
}
```

**Errors**

- `400` - Invalid request format or negative amount
- `401` - Missing or invalid token

**Rules**

- Updates the base income for current month
- Amount must be non-negative
- Creates budget if it doesn't exist
- Can be set to zero
- Can be updated multiple times per month

---

## Expense Endpoints

### POST /expenses

Add a new expense to the current month.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**

```json
{
  "title": "Rent",
  "amount": 1200
}
```

**Response** (201 Created)

```json
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

**Errors**

- `400` - Missing title or amount, or invalid amount (≤0)
- `401` - Missing or invalid token

**Rules**

- Amount must be positive (> 0)
- Expense is automatically added to current month
- Expense gets unique ID (MongoDB ObjectId)
- If budget doesn't exist, it's created automatically

---

### PUT /expenses/:expenseId

Update an existing expense.

**Headers**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameter**

- `expenseId` - MongoDB ObjectId of the expense (from response)

**Request**

```json
{
  "title": "Updated Rent",
  "amount": 1300
}
```

**Response** (200 OK)

```json
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Updated Rent",
      "amount": 1300,
      "createdAt": "2026-01-31T10:30:00Z"
    }
  ],
  "remaining": 3700
}
```

**Errors**

- `400` - Invalid request format or amount
- `401` - Missing or invalid token
- `500` - Expense not found or doesn't belong to user

**Rules**

- Amount must be positive (> 0)
- Only updates expenses in user's budgets
- Can only update current month's expenses
- `createdAt` timestamp is not updated

---

### DELETE /expenses/:expenseId

Delete an expense.

**Headers**

```
Authorization: Bearer <token>
```

**URL Parameter**

- `expenseId` - MongoDB ObjectId of the expense (from response)

**Response** (200 OK)

```json
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [],
  "remaining": 5000
}
```

**Errors**

- `401` - Missing or invalid token
- `500` - Expense not found or doesn't belong to user

**Rules**

- Only deletes user's own expenses
- Permanently removes the expense
- Budget remains if it's the last expense

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "descriptive error message"
}
```

**Common Error Messages**

- `"missing authorization header"` - No Authorization header provided
- `"invalid authorization header format"` - Wrong Bearer format
- `"invalid token"` - Token is expired or malformed
- `"invalid token claims"` - Token lacks userId claim
- `"invalid request format"` - JSON parsing failed
- `"email and password are required"` - Missing fields
- `"user with this email already exists"` - Duplicate email
- `"invalid email or password"` - Login failed
- `"invalid user ID"` - Malformed user ID
- `"invalid expense ID"` - Malformed expense ID
- `"title and amount (positive) are required"` - Missing or invalid expense fields
- `"amount must be non-negative"` - Negative base income
- `"expense not found or doesn't belong to user"` - Expense doesn't exist or unauthorized
- `"endpoint not found"` - 404 Not Found

---

## Authentication Details

### JWT Token Format

Tokens use HS256 signing and contain:

- `userId` - The authenticated user's ID
- `exp` - Expiration timestamp
- `iat` - Issued at timestamp

### Token Validation

- Tokens expire after 24 hours (configurable)
- Required in `Authorization: Bearer <token>` header
- Invalid or expired tokens return `401 Unauthorized`

### Security

- Passwords are hashed with bcrypt (cost 10)
- JWT secret should be strong and kept secure
- Change JWT_SECRET in production

---

## Rate Limiting

Currently no rate limiting is implemented. This should be added for production.

---

## Pagination

Pagination is not supported in this version. Consider adding for large datasets in future versions.

---

## Filtering & Sorting

Filtering and sorting are not supported in this version.

---

## Versioning

API version: 1.0.0
No versioning in URLs currently. Consider adding `/v1/` prefix for future versions.
