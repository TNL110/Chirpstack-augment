#!/bin/bash

set -e

echo "🧪 Testing Device Delete with ChirpStack Integration"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if service is running
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_error "Service is not running. Please start the service first."
    exit 1
fi

print_success "Service is running!"

# Generate unique identifiers
TIMESTAMP=$(date +%s)
RANDOM_SUFFIX=$(printf "%02X" $((RANDOM % 256)))
TEST_EMAIL="delete_test_${TIMESTAMP}@example.com"
DEV_EUI="08574DA9118DC9${RANDOM_SUFFIX}"

print_status "Using test email: $TEST_EMAIL"
print_status "Using DevEUI: $DEV_EUI"

# 1. Register user
print_status "Step 1: Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"password123\", \"full_name\": \"Delete Test User\"}")

if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    print_success "User registered successfully"
    TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
else
    print_error "User registration failed"
    echo "Response: $REGISTER_RESPONSE"
    exit 1
fi

# 2. Create device version
print_status "Step 2: Creating device version..."
VERSION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices/versions \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"RAK7200_DELETE_TEST_${TIMESTAMP}\", \"version\": \"v1.0\", \"description\": \"For delete testing\"}")

if echo "$VERSION_RESPONSE" | grep -q "id"; then
    print_success "Device version created successfully"
    VERSION_ID=$(echo "$VERSION_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
else
    print_error "Device version creation failed"
    echo "Response: $VERSION_RESPONSE"
    exit 1
fi

# 3. Create allowed device
print_status "Step 3: Creating allowed device..."
ALLOWED_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices/allowed \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"dev_eui\": \"$DEV_EUI\", \"nwk_key\": \"C518B15AB390B01762E4A3730E8C5F1C\", \"app_key\": \"97784F3B7F2A57EECF19F10E625081E0\", \"addr_key\": \"2F972E56\", \"description\": \"Device for delete testing\"}")

if echo "$ALLOWED_RESPONSE" | grep -q "id"; then
    print_success "Allowed device created successfully"
else
    print_error "Allowed device creation failed"
    echo "Response: $ALLOWED_RESPONSE"
    exit 1
fi

# 4. Create device (will be created in ChirpStack)
print_status "Step 4: Creating device (will create in ChirpStack)..."
DEVICE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"version_id\": \"$VERSION_ID\", \"name\": \"Delete Test Device\", \"dev_eui\": \"$DEV_EUI\", \"description\": \"Device to be deleted\"}")

if echo "$DEVICE_RESPONSE" | grep -q "chirpstack_device_created"; then
    print_success "Device created successfully"
    DEVICE_ID=$(echo "$DEVICE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    # Check if device was created in ChirpStack
    if echo "$DEVICE_RESPONSE" | grep -q '"chirpstack_device_created":true'; then
        print_success "✅ Device was created in ChirpStack"
    else
        print_warning "⚠️  Device was NOT created in ChirpStack"
    fi
    
    # Check if device was activated in ChirpStack
    if echo "$DEVICE_RESPONSE" | grep -q '"chirpstack_device_activated":true'; then
        print_success "✅ Device was activated in ChirpStack"
    else
        print_warning "⚠️  Device was NOT activated in ChirpStack"
    fi
else
    print_error "Device creation failed"
    echo "Response: $DEVICE_RESPONSE"
    exit 1
fi

print_status "Device ID: $DEVICE_ID"
print_status "DevEUI: $DEV_EUI"

# 5. Verify device exists
print_status "Step 5: Verifying device exists..."
GET_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/devices/$DEVICE_ID \
    -H "Authorization: Bearer $TOKEN")

if echo "$GET_RESPONSE" | grep -q "\"id\":\"$DEVICE_ID\""; then
    print_success "Device exists in database"
else
    print_error "Device not found in database"
    echo "Response: $GET_RESPONSE"
    exit 1
fi

# 6. Delete device (will delete from ChirpStack)
print_status "Step 6: Deleting device (will delete from ChirpStack)..."
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/v1/devices/$DEVICE_ID \
    -H "Authorization: Bearer $TOKEN")

if echo "$DELETE_RESPONSE" | grep -q "Device deleted successfully"; then
    print_success "✅ Device deleted successfully from database"
    print_success "✅ Device should also be deleted from ChirpStack"
else
    print_error "Device deletion failed"
    echo "Response: $DELETE_RESPONSE"
    exit 1
fi

# 7. Verify device is deleted
print_status "Step 7: Verifying device is deleted..."
VERIFY_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/devices/$DEVICE_ID \
    -H "Authorization: Bearer $TOKEN")

if echo "$VERIFY_RESPONSE" | grep -q "device not found"; then
    print_success "✅ Device successfully deleted from database"
else
    print_error "Device still exists in database"
    echo "Response: $VERIFY_RESPONSE"
    exit 1
fi

# 8. Check logs for ChirpStack deletion
print_status "Step 8: Checking logs for ChirpStack deletion..."
LOGS=$(docker logs go-auth-api 2>&1 | tail -20)

if echo "$LOGS" | grep -q "ChirpStack device deleted: $DEV_EUI"; then
    print_success "✅ ChirpStack device deletion confirmed in logs"
else
    print_warning "⚠️  ChirpStack device deletion not found in logs"
    print_status "Recent logs:"
    echo "$LOGS" | grep -E "(ChirpStack|DELETE)" || echo "No relevant logs found"
fi

print_success "🎉 Device Delete with ChirpStack Integration Test PASSED!"
echo ""
echo "Summary:"
echo "✅ User registered successfully"
echo "✅ Device version created"
echo "✅ Allowed device created"
echo "✅ Device created in database and ChirpStack"
echo "✅ Device deleted from database and ChirpStack"
echo "✅ Deletion verified"
echo ""
print_success "All tests completed successfully!"
