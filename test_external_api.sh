#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration - Change this to your server IP
SERVER_IP="192.168.0.93"
API_BASE="http://$SERVER_IP:8080"

echo -e "${YELLOW}=== Testing External API Access ===${NC}"
echo -e "${YELLOW}Server IP: $SERVER_IP${NC}"
echo -e "${YELLOW}API Base URL: $API_BASE${NC}"

# Test 1: Health Check
echo -e "\n${YELLOW}1. Testing Health Check...${NC}"
health_response=$(curl -s -w "%{http_code}" -o /tmp/health_external.json \
  -X GET "$API_BASE/health")

health_http_code="${health_response: -3}"

if [ "$health_http_code" = "200" ]; then
    echo -e "${GREEN}✓ Health check successful${NC}"
    cat /tmp/health_external.json
else
    echo -e "${RED}✗ Health check failed (HTTP $health_http_code)${NC}"
    cat /tmp/health_external.json
    exit 1
fi

# Test 2: User Registration
echo -e "\n${YELLOW}2. Testing User Registration...${NC}"

# Generate unique email for testing
TIMESTAMP=$(date +%s)
TEST_EMAIL="external_test_${TIMESTAMP}@example.com"
TEST_PASSWORD="password123"
TEST_NAME="External Test User $TIMESTAMP"

register_response=$(curl -s -w "%{http_code}" -o /tmp/register_external.json \
  -X POST "$API_BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"full_name\": \"$TEST_NAME\"
  }")

register_http_code="${register_response: -3}"

if [ "$register_http_code" = "201" ]; then
    echo -e "${GREEN}✓ User registration successful${NC}"
    echo "Response:"
    cat /tmp/register_external.json | jq '.'
    
    # Extract token for later use
    TOKEN=$(cat /tmp/register_external.json | jq -r '.token')
    echo -e "\n${GREEN}Token extracted for further tests${NC}"
else
    echo -e "${RED}✗ User registration failed (HTTP $register_http_code)${NC}"
    cat /tmp/register_external.json
    exit 1
fi

# Test 3: User Login
echo -e "\n${YELLOW}3. Testing User Login...${NC}"

login_response=$(curl -s -w "%{http_code}" -o /tmp/login_external.json \
  -X POST "$API_BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

login_http_code="${login_response: -3}"

if [ "$login_http_code" = "200" ]; then
    echo -e "${GREEN}✓ User login successful${NC}"
    echo "Response:"
    cat /tmp/login_external.json | jq '.'
else
    echo -e "${RED}✗ User login failed (HTTP $login_http_code)${NC}"
    cat /tmp/login_external.json
    exit 1
fi

# Test 4: Get User Profile (Protected Route)
echo -e "\n${YELLOW}4. Testing Protected Route (Get Profile)...${NC}"

profile_response=$(curl -s -w "%{http_code}" -o /tmp/profile_external.json \
  -X GET "$API_BASE/api/v1/user/profile" \
  -H "Authorization: Bearer $TOKEN")

profile_http_code="${profile_response: -3}"

if [ "$profile_http_code" = "200" ]; then
    echo -e "${GREEN}✓ Get profile successful${NC}"
    echo "Response:"
    cat /tmp/profile_external.json | jq '.'
else
    echo -e "${RED}✗ Get profile failed (HTTP $profile_http_code)${NC}"
    cat /tmp/profile_external.json
fi

# Test 5: Get All Users (Protected Route)
echo -e "\n${YELLOW}5. Testing Get All Users...${NC}"

users_response=$(curl -s -w "%{http_code}" -o /tmp/users_external.json \
  -X GET "$API_BASE/api/v1/users?page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")

users_http_code="${users_response: -3}"

if [ "$users_http_code" = "200" ]; then
    echo -e "${GREEN}✓ Get all users successful${NC}"
    echo "Response:"
    cat /tmp/users_external.json | jq '.'
else
    echo -e "${RED}✗ Get all users failed (HTTP $users_http_code)${NC}"
    cat /tmp/users_external.json
fi

# Test 6: Test without Authorization (should fail)
echo -e "\n${YELLOW}6. Testing Unauthorized Access...${NC}"

unauth_response=$(curl -s -w "%{http_code}" -o /tmp/unauth_external.json \
  -X GET "$API_BASE/api/v1/user/profile")

unauth_http_code="${unauth_response: -3}"

if [ "$unauth_http_code" = "401" ]; then
    echo -e "${GREEN}✓ Unauthorized access properly blocked${NC}"
    echo "Response:"
    cat /tmp/unauth_external.json | jq '.'
else
    echo -e "${RED}✗ Unauthorized access not properly blocked (HTTP $unauth_http_code)${NC}"
    cat /tmp/unauth_external.json
fi

echo -e "\n${YELLOW}=== External API Test Summary ===${NC}"
echo -e "Server IP: $SERVER_IP"
echo -e "Test User: $TEST_EMAIL"
echo -e "All tests completed. Check results above."

# Cleanup
rm -f /tmp/*_external.json

echo -e "\n${GREEN}=== API Endpoints Available Externally ===${NC}"
echo -e "Health Check: ${GREEN}GET $API_BASE/health${NC}"
echo -e "Register: ${GREEN}POST $API_BASE/api/v1/auth/register${NC}"
echo -e "Login: ${GREEN}POST $API_BASE/api/v1/auth/login${NC}"
echo -e "Profile: ${GREEN}GET $API_BASE/api/v1/user/profile${NC} (requires token)"
echo -e "Users: ${GREEN}GET $API_BASE/api/v1/users${NC} (requires token)"
echo -e "Search: ${GREEN}GET $API_BASE/api/v1/users/search${NC} (requires token)"
