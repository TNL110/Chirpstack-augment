#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"

echo -e "${YELLOW}=== Testing Go Auth API ===${NC}"

# Test 1: Health Check
echo -e "\n${YELLOW}1. Testing Health Check...${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/health_response.json "$API_BASE/health")
http_code="${response: -3}"

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    cat /tmp/health_response.json
else
    echo -e "${RED}✗ Health check failed (HTTP $http_code)${NC}"
    cat /tmp/health_response.json
fi

# Test 2: Register User
echo -e "\n${YELLOW}2. Testing User Registration...${NC}"
register_response=$(curl -s -w "%{http_code}" -o /tmp/register_response.json \
  -X POST "$API_BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }')

register_http_code="${register_response: -3}"

if [ "$register_http_code" = "201" ]; then
    echo -e "${GREEN}✓ User registration successful${NC}"
    cat /tmp/register_response.json
    # Extract token for later use
    TOKEN=$(cat /tmp/register_response.json | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    echo -e "${RED}✗ User registration failed (HTTP $register_http_code)${NC}"
    cat /tmp/register_response.json
fi

# Test 3: Login User
echo -e "\n${YELLOW}3. Testing User Login...${NC}"
login_response=$(curl -s -w "%{http_code}" -o /tmp/login_response.json \
  -X POST "$API_BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

login_http_code="${login_response: -3}"

if [ "$login_http_code" = "200" ]; then
    echo -e "${GREEN}✓ User login successful${NC}"
    cat /tmp/login_response.json
    # Extract token for later use (in case registration failed)
    if [ -z "$TOKEN" ]; then
        TOKEN=$(cat /tmp/login_response.json | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    fi
else
    echo -e "${RED}✗ User login failed (HTTP $login_http_code)${NC}"
    cat /tmp/login_response.json
fi

# Test 4: Access Protected Route
if [ ! -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}4. Testing Protected Route (Profile)...${NC}"
    profile_response=$(curl -s -w "%{http_code}" -o /tmp/profile_response.json \
      -X GET "$API_BASE/api/v1/user/profile" \
      -H "Authorization: Bearer $TOKEN")

    profile_http_code="${profile_response: -3}"

    if [ "$profile_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Protected route access successful${NC}"
        cat /tmp/profile_response.json
    else
        echo -e "${RED}✗ Protected route access failed (HTTP $profile_http_code)${NC}"
        cat /tmp/profile_response.json
    fi
else
    echo -e "\n${RED}4. Skipping protected route test - no token available${NC}"
fi

# Test 5: Access Protected Route without Token
echo -e "\n${YELLOW}5. Testing Protected Route without Token...${NC}"
no_auth_response=$(curl -s -w "%{http_code}" -o /tmp/no_auth_response.json \
  -X GET "$API_BASE/api/v1/user/profile")

no_auth_http_code="${no_auth_response: -3}"

if [ "$no_auth_http_code" = "401" ]; then
    echo -e "${GREEN}✓ Unauthorized access properly blocked${NC}"
    cat /tmp/no_auth_response.json
else
    echo -e "${RED}✗ Unauthorized access not properly blocked (HTTP $no_auth_http_code)${NC}"
    cat /tmp/no_auth_response.json
fi

# Test 6: Try to register duplicate user
echo -e "\n${YELLOW}6. Testing Duplicate User Registration...${NC}"
duplicate_response=$(curl -s -w "%{http_code}" -o /tmp/duplicate_response.json \
  -X POST "$API_BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User Duplicate"
  }')

duplicate_http_code="${duplicate_response: -3}"

if [ "$duplicate_http_code" = "400" ]; then
    echo -e "${GREEN}✓ Duplicate user registration properly blocked${NC}"
    cat /tmp/duplicate_response.json
else
    echo -e "${RED}✗ Duplicate user registration not properly blocked (HTTP $duplicate_http_code)${NC}"
    cat /tmp/duplicate_response.json
fi

# Test 7: Get All Users
if [ ! -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}7. Testing Get All Users...${NC}"
    users_response=$(curl -s -w "%{http_code}" -o /tmp/users_response.json \
      -X GET "$API_BASE/api/v1/users?page=1&page_size=5" \
      -H "Authorization: Bearer $TOKEN")

    users_http_code="${users_response: -3}"

    if [ "$users_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Get all users successful${NC}"
        cat /tmp/users_response.json
    else
        echo -e "${RED}✗ Get all users failed (HTTP $users_http_code)${NC}"
        cat /tmp/users_response.json
    fi
else
    echo -e "\n${RED}7. Skipping get all users test - no token available${NC}"
fi

# Test 8: Get User by ID
if [ ! -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}8. Testing Get User by ID...${NC}"
    user_by_id_response=$(curl -s -w "%{http_code}" -o /tmp/user_by_id_response.json \
      -X GET "$API_BASE/api/v1/users/1" \
      -H "Authorization: Bearer $TOKEN")

    user_by_id_http_code="${user_by_id_response: -3}"

    if [ "$user_by_id_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Get user by ID successful${NC}"
        cat /tmp/user_by_id_response.json
    else
        echo -e "${RED}✗ Get user by ID failed (HTTP $user_by_id_http_code)${NC}"
        cat /tmp/user_by_id_response.json
    fi
else
    echo -e "\n${RED}8. Skipping get user by ID test - no token available${NC}"
fi

# Test 9: Update User
if [ ! -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}9. Testing Update User...${NC}"
    update_response=$(curl -s -w "%{http_code}" -o /tmp/update_response.json \
      -X PUT "$API_BASE/api/v1/users/1" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "full_name": "Updated Test User"
      }')

    update_http_code="${update_response: -3}"

    if [ "$update_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Update user successful${NC}"
        cat /tmp/update_response.json
    else
        echo -e "${RED}✗ Update user failed (HTTP $update_http_code)${NC}"
        cat /tmp/update_response.json
    fi
else
    echo -e "\n${RED}9. Skipping update user test - no token available${NC}"
fi

# Test 10: Search Users
if [ ! -z "$TOKEN" ]; then
    echo -e "\n${YELLOW}10. Testing Search Users...${NC}"
    search_response=$(curl -s -w "%{http_code}" -o /tmp/search_response.json \
      -X GET "$API_BASE/api/v1/users/search?q=test&page=1&page_size=5" \
      -H "Authorization: Bearer $TOKEN")

    search_http_code="${search_response: -3}"

    if [ "$search_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Search users successful${NC}"
        cat /tmp/search_response.json
    else
        echo -e "${RED}✗ Search users failed (HTTP $search_http_code)${NC}"
        cat /tmp/search_response.json
    fi
else
    echo -e "\n${RED}10. Skipping search users test - no token available${NC}"
fi

echo -e "\n${YELLOW}=== Test Summary ===${NC}"
echo "All tests completed. Check the results above."

# Cleanup
rm -f /tmp/*_response.json
