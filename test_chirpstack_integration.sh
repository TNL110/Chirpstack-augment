#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"
CHIRPSTACK_BASE="http://192.168.0.21:8090"
CHIRPSTACK_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJjaGlycHN0YWNrIiwiaXNzIjoiY2hpcnBzdGFjayIsInN1YiI6Ijg5MjIwOTc5LTI2NWYtNGYyOC1hYmIxLTBiMDNlOWQwMDA5MyIsInR5cCI6ImtleSJ9.BWpMBTWTFPNFPjrW-_ELN2-bK4VTKMeXQK6HkaijJaE"

echo -e "${YELLOW}=== Testing ChirpStack Integration ===${NC}"

# Test 1: Check ChirpStack connectivity
echo -e "\n${YELLOW}1. Testing ChirpStack API connectivity...${NC}"
chirpstack_response=$(curl -s -w "%{http_code}" -o /tmp/chirpstack_test.json \
  -X GET "$CHIRPSTACK_BASE/api/tenants" \
  -H "Authorization: Bearer $CHIRPSTACK_TOKEN")

chirpstack_http_code="${chirpstack_response: -3}"

if [ "$chirpstack_http_code" = "200" ]; then
    echo -e "${GREEN}✓ ChirpStack API is accessible${NC}"
else
    echo -e "${RED}✗ ChirpStack API is not accessible (HTTP $chirpstack_http_code)${NC}"
    cat /tmp/chirpstack_test.json
    echo -e "${RED}Please check ChirpStack server and token${NC}"
    exit 1
fi

# Test 2: Register a new user with ChirpStack integration
echo -e "\n${YELLOW}2. Testing User Registration with ChirpStack Integration...${NC}"

# Generate unique email for testing
TIMESTAMP=$(date +%s)
TEST_EMAIL="chirpstack_test_${TIMESTAMP}@example.com"
TEST_USERNAME="chirpstack_test_${TIMESTAMP}"

register_response=$(curl -s -w "%{http_code}" -o /tmp/register_chirpstack.json \
  -X POST "$API_BASE/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"password123\",
    \"full_name\": \"ChirpStack Test User\"
  }")

register_http_code="${register_response: -3}"

if [ "$register_http_code" = "201" ]; then
    echo -e "${GREEN}✓ User registration successful${NC}"
    cat /tmp/register_chirpstack.json
    
    # Extract token for later use
    TOKEN=$(cat /tmp/register_chirpstack.json | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo -e "\n${GREEN}Token extracted: ${TOKEN:0:50}...${NC}"
else
    echo -e "${RED}✗ User registration failed (HTTP $register_http_code)${NC}"
    cat /tmp/register_chirpstack.json
    exit 1
fi

# Test 3: Check if ChirpStack resources were created
echo -e "\n${YELLOW}3. Checking ChirpStack Resources Creation...${NC}"

# Wait a moment for resources to be created
sleep 2

# Check tenants
echo -e "\n${YELLOW}3.1. Checking Tenants...${NC}"
tenants_response=$(curl -s -w "%{http_code}" -o /tmp/tenants_check.json \
  -X GET "$CHIRPSTACK_BASE/api/tenants?limit=100" \
  -H "Authorization: Bearer $CHIRPSTACK_TOKEN")

tenants_http_code="${tenants_response: -3}"

if [ "$tenants_http_code" = "200" ]; then
    echo -e "${GREEN}✓ Successfully retrieved tenants${NC}"
    
    # Check if our tenant exists
    if grep -q "$TEST_EMAIL" /tmp/tenants_check.json; then
        echo -e "${GREEN}✓ Tenant for user $TEST_EMAIL found${NC}"
        # Extract tenant ID for the specific user
        TENANT_ID=$(cat /tmp/tenants_check.json | jq -r ".result[] | select(.name == \"$TEST_EMAIL\") | .id")
        echo -e "${GREEN}Tenant ID: $TENANT_ID${NC}"
    else
        echo -e "${RED}✗ Tenant for user $TEST_EMAIL not found${NC}"
        echo "Available tenants:"
        cat /tmp/tenants_check.json | jq '.'
    fi
else
    echo -e "${RED}✗ Failed to retrieve tenants (HTTP $tenants_http_code)${NC}"
    cat /tmp/tenants_check.json
fi

# Test 4: Check Applications
if [ ! -z "$TENANT_ID" ]; then
    echo -e "\n${YELLOW}3.2. Checking Applications...${NC}"
    apps_response=$(curl -s -w "%{http_code}" -o /tmp/apps_check.json \
      -X GET "$CHIRPSTACK_BASE/api/applications?tenantId=$TENANT_ID&limit=100" \
      -H "Authorization: Bearer $CHIRPSTACK_TOKEN")

    apps_http_code="${apps_response: -3}"

    if [ "$apps_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Successfully retrieved applications${NC}"
        
        # Check if Lnode application exists
        if grep -q "Lnode" /tmp/apps_check.json; then
            echo -e "${GREEN}✓ Lnode application found${NC}"
            APP_ID=$(cat /tmp/apps_check.json | jq -r '.result[] | select(.name == "Lnode") | .id')
            echo -e "${GREEN}Application ID: $APP_ID${NC}"
        else
            echo -e "${RED}✗ Lnode application not found${NC}"
            echo "Available applications:"
            cat /tmp/apps_check.json | jq '.'
        fi
    else
        echo -e "${RED}✗ Failed to retrieve applications (HTTP $apps_http_code)${NC}"
        cat /tmp/apps_check.json
    fi
fi

# Test 5: Check Device Profiles
if [ ! -z "$TENANT_ID" ]; then
    echo -e "\n${YELLOW}3.3. Checking Device Profiles...${NC}"
    profiles_response=$(curl -s -w "%{http_code}" -o /tmp/profiles_check.json \
      -X GET "$CHIRPSTACK_BASE/api/device-profiles?tenantId=$TENANT_ID&limit=100" \
      -H "Authorization: Bearer $CHIRPSTACK_TOKEN")

    profiles_http_code="${profiles_response: -3}"

    if [ "$profiles_http_code" = "200" ]; then
        echo -e "${GREEN}✓ Successfully retrieved device profiles${NC}"
        
        # Check if RAK_ABP profile exists
        if grep -q "RAK_ABP" /tmp/profiles_check.json; then
            echo -e "${GREEN}✓ RAK_ABP device profile found${NC}"
            PROFILE_ID=$(cat /tmp/profiles_check.json | jq -r '.result[] | select(.name == "RAK_ABP") | .id')
            echo -e "${GREEN}Device Profile ID: $PROFILE_ID${NC}"
        else
            echo -e "${RED}✗ RAK_ABP device profile not found${NC}"
            echo "Available device profiles:"
            cat /tmp/profiles_check.json | jq '.'
        fi
    else
        echo -e "${RED}✗ Failed to retrieve device profiles (HTTP $profiles_http_code)${NC}"
        cat /tmp/profiles_check.json
    fi
fi

echo -e "\n${YELLOW}=== ChirpStack Integration Test Summary ===${NC}"
echo -e "Test completed for user: $TEST_EMAIL"
echo -e "Check the results above for any issues."

# Cleanup
rm -f /tmp/*_check.json /tmp/register_chirpstack.json /tmp/chirpstack_test.json
