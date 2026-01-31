#!/bin/bash
# Finance App API - Example Requests
# Make sure the server is running on localhost:3000

BASE_URL="http://localhost:3000"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Finance App API Testing ===${NC}\n"

# 1. Sign Up
echo -e "${GREEN}1. Signing up a new user...${NC}"
SIGNUP_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }')

echo "$SIGNUP_RESPONSE" | jq .
TOKEN=$(echo "$SIGNUP_RESPONSE" | jq -r '.token')
echo -e "Token: ${YELLOW}${TOKEN:0:30}...${NC}\n"

# 2. Log In
echo -e "${GREEN}2. Logging in with the same credentials...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }')

echo "$LOGIN_RESPONSE" | jq .
LOGIN_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
echo -e "Token: ${YELLOW}${LOGIN_TOKEN:0:30}...${NC}\n"

# 3. Get Current Budget
echo -e "${GREEN}3. Getting current month's budget...${NC}"
curl -s -X GET "${BASE_URL}/budget/current" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 4. Set Base Income
echo -e "${GREEN}4. Setting base income for current month...${NC}"
curl -s -X POST "${BASE_URL}/budget/base-income" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000
  }' | jq .
echo ""

# 5. Add Expense 1
echo -e "${GREEN}5. Adding first expense (Rent)...${NC}"
EXPENSE_RESPONSE=$(curl -s -X POST "${BASE_URL}/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Rent",
    "amount": 1200
  }')

echo "$EXPENSE_RESPONSE" | jq .
EXPENSE_ID=$(echo "$EXPENSE_RESPONSE" | jq -r '.expenses[0].id')
echo -e "Expense ID: ${YELLOW}${EXPENSE_ID}${NC}\n"

# 6. Add Expense 2
echo -e "${GREEN}6. Adding second expense (Groceries)...${NC}"
curl -s -X POST "${BASE_URL}/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Groceries",
    "amount": 300
  }' | jq .
echo ""

# 7. Add Expense 3
echo -e "${GREEN}7. Adding third expense (Utilities)...${NC}"
curl -s -X POST "${BASE_URL}/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Utilities",
    "amount": 150
  }' | jq .
echo ""

# 8. Get Current Budget (with expenses)
echo -e "${GREEN}8. Getting current budget with all expenses...${NC}"
BUDGET=$(curl -s -X GET "${BASE_URL}/budget/current" \
  -H "Authorization: Bearer $TOKEN")

echo "$BUDGET" | jq .
echo ""

# 9. Update Expense
echo -e "${GREEN}9. Updating the Rent expense amount to 1300...${NC}"
curl -s -X PUT "${BASE_URL}/expenses/${EXPENSE_ID}" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Rent (Updated)",
    "amount": 1300
  }' | jq .
echo ""

# 10. Get Specific Month Budget
echo -e "${GREEN}10. Getting budget for December 2025...${NC}"
curl -s -X GET "${BASE_URL}/budget?year=2025&month=12" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 11. Delete Expense
echo -e "${GREEN}11. Deleting the Utilities expense...${NC}"
curl -s -X DELETE "${BASE_URL}/expenses/${EXPENSE_ID}" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 12. Try unauthorized request (missing token)
echo -e "${GREEN}12. Testing unauthorized access (missing token)...${NC}"
curl -s -X GET "${BASE_URL}/budget/current" | jq .
echo ""

# 13. Test invalid token
echo -e "${GREEN}13. Testing with invalid token...${NC}"
curl -s -X GET "${BASE_URL}/budget/current" \
  -H "Authorization: Bearer invalid_token_123" | jq .
echo ""

echo -e "${GREEN}=== Testing Complete ===${NC}"
