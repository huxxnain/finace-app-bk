# Finance App Backend - Project Index

Welcome to the Personal Monthly Expense Tracker backend! This document helps you navigate all the files and get started quickly.

## 📚 Documentation (Start Here!)

1. **[README.md](README.md)** - Project overview, setup, and features
   - Best for: Understanding what the project does

2. **[QUICKSTART.md](QUICKSTART.md)** - Get running in 5 minutes
   - Best for: Quick setup with Docker or local Go
   - Contains: Step-by-step instructions

3. **[API.md](API.md)** - Complete API reference
   - Best for: Understanding all endpoints
   - Contains: Request/response examples, error codes

4. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Technical deep dive
   - Best for: Understanding code structure
   - Contains: Architecture diagrams, data flow, security

5. **[DEPLOYMENT.md](DEPLOYMENT.md)** - Production deployment
   - Best for: Deploying to servers
   - Contains: Docker, Fly.io, Render, VPS instructions

6. **[PROJECT_COMPLETION.md](PROJECT_COMPLETION.md)** - Completion summary
   - Best for: Overview of what was delivered
   - Contains: Features, requirements compliance, success criteria

---

## 📁 Project Structure

### Configuration Files

```
.env.example              ← Environment variables template
.gitignore               ← Git ignore rules
go.mod                   ← Go module dependencies
Makefile                 ← Build commands
```

### Source Code

```
cmd/
  └── server/
      └── main.go        ← Application entry point (Fiber server setup)

internal/
  ├── auth/              ← Authentication (JWT, middleware)
  │   ├── jwt.go
  │   └── middleware.go
  ├── config/            ← Configuration management
  │   └── config.go
  ├── db/                ← Database connection
  │   └── mongo.go
  ├── handlers/          ← HTTP handlers
  │   ├── auth_handler.go
  │   ├── budget_handler.go
  │   └── expense_handler.go
  ├── models/            ← Data structures
  │   └── models.go
  ├── services/          ← Business logic
  │   ├── user_service.go
  │   └── budget_service.go
  └── utils/             ← Helper functions
      ├── password.go
      └── time.go
```

### Docker & DevOps

```
Dockerfile               ← Multi-stage Docker build
docker-compose.yml      ← Docker Compose setup (app + MongoDB)
```

### Testing & Examples

```
examples.sh                                    ← 13 comprehensive curl examples
test-simple.sh                                 ← Quick test script
Finance-App-API.postman_collection.json       ← Postman collection import
```

---

## 🚀 Quick Navigation

### I want to...

**Get started quickly (5 min)**
→ Read [QUICKSTART.md](QUICKSTART.md)

**Understand how the API works**
→ Read [API.md](API.md)

**Understand the code architecture**
→ Read [ARCHITECTURE.md](ARCHITECTURE.md)

**Deploy to production**
→ Read [DEPLOYMENT.md](DEPLOYMENT.md)

**Test the API**
→ Run `bash examples.sh` or import Postman collection

**Understand requirements**
→ Read [PROJECT_COMPLETION.md](PROJECT_COMPLETION.md)

**See all features**
→ Read [README.md](README.md)

---

## 📖 Documentation Quick Reference

### Setup & Getting Started

- [README.md](README.md) - Full setup instructions
- [QUICKSTART.md](QUICKSTART.md) - Fast setup (5 min)
- `.env.example` - Environment template

### API Usage

- [API.md](API.md) - All endpoints and examples
- [Finance-App-API.postman_collection.json](Finance-App-API.postman_collection.json) - Import into Postman
- `examples.sh` - curl examples

### Code & Architecture

- [ARCHITECTURE.md](ARCHITECTURE.md) - Code structure
- Source files in `cmd/` and `internal/`
- Comments within source files

### Deployment

- [DEPLOYMENT.md](DEPLOYMENT.md) - All deployment options
- `Dockerfile` - Docker build
- `docker-compose.yml` - Full stack with MongoDB

### Project Info

- [PROJECT_COMPLETION.md](PROJECT_COMPLETION.md) - Delivery summary
- `Makefile` - Build commands
- `go.mod` - Dependencies

---

## 🎯 Key Files by Purpose

### Authentication

- `internal/auth/jwt.go` - JWT token generation/validation
- `internal/auth/middleware.go` - Bearer token middleware
- `internal/handlers/auth_handler.go` - /auth/signup, /auth/login endpoints
- `internal/services/user_service.go` - User registration and login logic
- `internal/utils/password.go` - Bcrypt password hashing

### Budget Management

- `internal/handlers/budget_handler.go` - Budget endpoints
- `internal/services/budget_service.go` - Budget business logic
- `internal/models/models.go` - MonthlyBudget data model

### Expense Management

- `internal/handlers/expense_handler.go` - Expense endpoints
- `internal/services/budget_service.go` - Expense operations
- `internal/models/models.go` - Expense data model

### Database

- `internal/db/mongo.go` - MongoDB connection
- `internal/config/config.go` - Database configuration

### Server & Routing

- `cmd/server/main.go` - Fiber app setup, routes, middleware

---

## 🧪 Testing Quick Start

### Option 1: Run comprehensive tests

```bash
bash examples.sh
```

### Option 2: Run quick test

```bash
bash test-simple.sh
```

### Option 3: Manual testing with curl

```bash
# 1. Sign up
TOKEN=$(curl -X POST http://localhost:3000/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}' \
  | jq -r '.token')

# 2. Set base income
curl -X POST http://localhost:3000/budget/base-income \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":5000}'

# 3. Add expense
curl -X POST http://localhost:3000/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Rent","amount":1200}'

# 4. Check balance
curl -X GET http://localhost:3000/budget/current \
  -H "Authorization: Bearer $TOKEN"
```

### Option 4: Use Postman

1. Import `Finance-App-API.postman_collection.json`
2. Set `base_url` to `http://localhost:3000`
3. Set `token` from signup response
4. Run requests

---

## 📊 Project Statistics

- **Total Files**: 28
- **Go Source Files**: 13
- **Documentation Files**: 6
- **Configuration Files**: 7
- **Lines of Code**: ~400 (source) + ~3000 (docs)
- **Test Scripts**: 2
- **Docker Setup**: 2 files

---

## ✅ Core Features

### Authentication

✅ Sign up with email/password
✅ Login with email/password
✅ JWT tokens (24hr expiry)
✅ Bcrypt password hashing
✅ Bearer token authentication

### Budget Management

✅ Get current month budget
✅ Get specific month budget
✅ Automatic budget creation
✅ Set base income
✅ Calculate remaining balance

### Expense Management

✅ Add expenses
✅ Update expenses
✅ Delete expenses
✅ Track expense titles and amounts
✅ Per-user expense isolation

### Security

✅ User data isolation
✅ JWT authentication
✅ Password hashing
✅ Input validation
✅ Error handling

---

## 🔧 Common Commands

### Build & Run

```bash
go run cmd/server/main.go        # Run directly
go build -o finance-app cmd/...  # Build binary
make run                         # Using Makefile
```

### Docker

```bash
docker-compose up -d             # Start with Docker Compose
docker build -t finance-app .    # Build image
```

### Testing

```bash
bash examples.sh                  # Full test suite
bash test-simple.sh              # Quick test
go test ./...                     # Go tests
```

### Maintenance

```bash
go mod tidy                       # Clean dependencies
go fmt ./...                      # Format code
make clean                        # Clean build artifacts
```

---

## 📞 Support & Help

### For Setup Issues

→ Read [QUICKSTART.md](QUICKSTART.md) - Troubleshooting section

### For API Questions

→ Read [API.md](API.md) - All endpoints documented

### For Code Understanding

→ Read [ARCHITECTURE.md](ARCHITECTURE.md) - Design and structure

### For Deployment

→ Read [DEPLOYMENT.md](DEPLOYMENT.md) - All platforms covered

### For Testing

→ Run `bash examples.sh` and `bash test-simple.sh`

---

## 📋 Checklist for First-Time Setup

- [ ] Read [QUICKSTART.md](QUICKSTART.md)
- [ ] Copy `.env.example` to `.env`
- [ ] Start MongoDB (Docker or local)
- [ ] Run `go mod download`
- [ ] Run `go run cmd/server/main.go`
- [ ] Test with `bash examples.sh`
- [ ] Read [API.md](API.md) for API details
- [ ] Import Postman collection for manual testing

---

## 🎓 Learning Path

1. **Start**: [QUICKSTART.md](QUICKSTART.md) - Get it running
2. **Test**: `bash examples.sh` - See it work
3. **Learn**: [API.md](API.md) - Understand endpoints
4. **Understand**: [ARCHITECTURE.md](ARCHITECTURE.md) - Deep dive
5. **Deploy**: [DEPLOYMENT.md](DEPLOYMENT.md) - Production setup

---

## 📝 Notes

- All endpoints return JSON
- All protected endpoints need Bearer token
- MongoDB required (or use Docker Compose)
- Environment variables in `.env` file
- Change JWT_SECRET in production
- Passwords hashed with bcrypt

---

## 🎉 You're All Set!

Your Finance App backend is ready to go. Start with [QUICKSTART.md](QUICKSTART.md) and enjoy building! 🚀

---

**Last Updated**: January 31, 2026
**Version**: 1.0.0
**Status**: Production Ready ✅
