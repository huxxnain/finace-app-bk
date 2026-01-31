# Quick Start Guide

Get your Finance App backend running in 5 minutes!

## Prerequisites

- Go 1.21+ OR Docker & Docker Compose
- MongoDB running locally OR Docker

## Option 1: Run Locally (Go)

### Step 1: Setup

```bash
# Clone/navigate to the project
cd finance-app

# Install dependencies
go mod download
go mod tidy

# Copy environment file
cp .env.example .env
```

### Step 2: Start MongoDB

```bash
# Using Docker (recommended)
docker run -d -p 27017:27017 --name finance-mongo mongo:latest

# OR using local MongoDB
mongod
```

### Step 3: Run the Server

```bash
go run cmd/server/main.go
```

You should see:

```
Server starting on port 3000
```

### Step 4: Test the API

```bash
# In another terminal, sign up
curl -X POST http://localhost:3000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

✅ You should get a token back!

---

## Option 2: Run with Docker Compose

### Step 1: Setup

```bash
cd finance-app
cp .env.example .env
```

### Step 2: Start Everything

```bash
docker-compose up -d
```

This starts:

- MongoDB on `localhost:27017`
- Finance App on `localhost:3000`

### Step 3: Test the API

```bash
curl -X POST http://localhost:3000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

✅ You should get a token back!

### Stop Everything

```bash
docker-compose down
```

---

## Using the API

### 1. Sign Up

```bash
TOKEN=$(curl -X POST http://localhost:3000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}' \
  | jq -r '.token')

echo $TOKEN  # Save this token
```

### 2. Set Base Income

```bash
curl -X POST http://localhost:3000/budget/base-income \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":5000}'
```

### 3. Add Expenses

```bash
curl -X POST http://localhost:3000/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Rent","amount":1200}'

curl -X POST http://localhost:3000/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Groceries","amount":300}'
```

### 4. Check Balance

```bash
curl -X GET http://localhost:3000/budget/current \
  -H "Authorization: Bearer $TOKEN"
```

Should show:

```json
{
  "year": 2026,
  "month": 1,
  "baseIncome": 5000,
  "expenses": [
    { "id": "...", "title": "Rent", "amount": 1200, "createdAt": "..." },
    { "id": "...", "title": "Groceries", "amount": 300, "createdAt": "..." }
  ],
  "remaining": 3500
}
```

---

## Useful Commands

### Run Tests

```bash
go test ./... -v
```

### Build for Production

```bash
go build -o finance-app cmd/server/main.go
```

### Format Code

```bash
go fmt ./...
```

### View Logs

```bash
# If using docker-compose
docker-compose logs -f finance-app

# If using docker run
docker logs -f finance-app
```

---

## Troubleshooting

### MongoDB Connection Error

- Ensure MongoDB is running: `docker ps | grep mongo`
- Check `MONGODB_URI` in `.env`
- Default: `mongodb://localhost:27017`

### Port Already in Use

- Change `PORT` in `.env` (default: 3000)
- Or kill the process using port 3000

### Invalid Token Error

- Token may have expired (24 hours by default)
- Get a new token by logging in again
- Or change `JWT_EXPIRY_HOURS` in `.env`

### CORS Issues

- Frontend requests blocked? Add CORS middleware to `cmd/server/main.go`
- See full documentation in `API.md`

---

## Next Steps

1. Read the [README.md](README.md) for full documentation
2. Check [API.md](API.md) for detailed API reference
3. Import [Finance-App-API.postman_collection.json](Finance-App-API.postman_collection.json) into Postman
4. Run `bash examples.sh` for comprehensive API testing

---

## Environment Variables

Key configurations in `.env`:

```
MONGODB_URI=mongodb://localhost:27017   # MongoDB connection
DATABASE_NAME=finance_app                # Database name
JWT_SECRET=your-super-secret-key         # Change this!
JWT_EXPIRY_HOURS=24                      # Token expiration
PORT=3000                                # Server port
```

---

## Support

Check [API.md](API.md) for complete API documentation
Check [README.md](README.md) for architecture details

Happy budgeting! 🎯
