#!/bin/bash
# Minimal Testing Script - Test only core functionality
# Make sure the server is running on localhost:3000

BASE_URL="http://localhost:3000"

echo "Testing Finance App API..."
echo ""

# 1. Sign Up
echo "1. Sign Up"
RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}')
echo "$RESPONSE" | jq .
TOKEN=$(echo "$RESPONSE" | jq -r '.token')
echo ""

# 2. Set Base Income
echo "2. Set Base Income to 5000"
curl -s -X POST "${BASE_URL}/budget/base-income" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":5000}' | jq .
echo ""

# 3. Add Expense
echo "3. Add Expense"
curl -s -X POST "${BASE_URL}/expenses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Rent","amount":1200}' | jq .
echo ""

# 4. Get Current Budget
echo "4. Get Current Budget"
curl -s -X GET "${BASE_URL}/budget/current" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "Testing complete!"
