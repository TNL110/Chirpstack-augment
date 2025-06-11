#!/bin/bash

set -e

echo "🧪 Starting Go Auth API Test Suite"
echo "=================================="

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

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Function to wait for service to be ready
wait_for_service() {
    print_status "Waiting for service to be ready..."
    max_attempts=30
    attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_success "Service is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_error "Service failed to start within 60 seconds"
    return 1
}

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    # Install test dependencies
    go mod tidy
    
    # Run unit tests with coverage
    go test -v -race -coverprofile=coverage.out ./tests/... -run "Test.*" -short
    
    if [ $? -eq 0 ]; then
        print_success "Unit tests passed!"
        
        # Generate coverage report
        go tool cover -html=coverage.out -o coverage.html
        print_status "Coverage report generated: coverage.html"
        
        # Show coverage summary
        go tool cover -func=coverage.out | tail -1
    else
        print_error "Unit tests failed!"
        return 1
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    
    # Start services
    print_status "Starting services..."
    docker-compose down -v > /dev/null 2>&1 || true
    docker-compose up --build -d
    
    # Wait for service to be ready
    if ! wait_for_service; then
        print_error "Failed to start services for integration tests"
        docker-compose logs
        docker-compose down
        return 1
    fi
    
    # Run integration tests
    go test -v ./tests/... -run "TestIntegration.*" -timeout 10m
    
    if [ $? -eq 0 ]; then
        print_success "Integration tests passed!"
    else
        print_error "Integration tests failed!"
        print_status "Service logs:"
        docker-compose logs app
        docker-compose down
        return 1
    fi
    
    # Clean up
    docker-compose down > /dev/null 2>&1
}

# Function to run API tests
run_api_tests() {
    print_status "Running API tests..."
    
    # Start services if not running
    if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_status "Starting services for API tests..."
        docker-compose up --build -d
        wait_for_service
    fi
    
    # Test user registration
    print_status "Testing user registration..."
    TIMESTAMP=$(date +%s)
    TEST_EMAIL="api_test_${TIMESTAMP}@example.com"
    REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"password123\", \"full_name\": \"API Test User\"}")
    
    if echo "$REGISTER_RESPONSE" | grep -q "token"; then
        print_success "User registration test passed"
        TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    else
        print_error "User registration test failed"
        echo "Response: $REGISTER_RESPONSE"
        return 1
    fi
    
    # Test user login
    print_status "Testing user login..."
    LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"password123\"}")
    
    if echo "$LOGIN_RESPONSE" | grep -q "token"; then
        print_success "User login test passed"
    else
        print_error "User login test failed"
        echo "Response: $LOGIN_RESPONSE"
        return 1
    fi
    
    # Test protected route
    print_status "Testing protected route..."
    PROFILE_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/user/profile \
        -H "Authorization: Bearer $TOKEN")
    
    if echo "$PROFILE_RESPONSE" | grep -q "email"; then
        print_success "Protected route test passed"
    else
        print_error "Protected route test failed"
        echo "Response: $PROFILE_RESPONSE"
        return 1
    fi
    
    # Test device version creation
    print_status "Testing device version creation..."
    TIMESTAMP=$(date +%s)
    VERSION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices/versions \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"name\": \"RAK7200_$TIMESTAMP\", \"version\": \"v1.0_$TIMESTAMP\", \"description\": \"Test version\"}")
    
    if echo "$VERSION_RESPONSE" | grep -q "id"; then
        print_success "Device version creation test passed"
        VERSION_ID=$(echo "$VERSION_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    else
        print_error "Device version creation test failed"
        echo "Response: $VERSION_RESPONSE"
        return 1
    fi
    
    # Test allowed device creation
    print_status "Testing allowed device creation..."
    DEV_EUI="C5EABC521E8304${TIMESTAMP:0:2}"
    ALLOWED_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices/allowed \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"dev_eui\": \"$DEV_EUI\", \"nwk_key\": \"C518B15AB390B01762E4A3730E8C5F1C\", \"app_key\": \"97784F3B7F2A57EECF19F10E625081E0\", \"addr_key\": \"2F972E56\", \"description\": \"Test device\"}")
    
    if echo "$ALLOWED_RESPONSE" | grep -q "id"; then
        print_success "Allowed device creation test passed"
    else
        print_error "Allowed device creation test failed"
        echo "Response: $ALLOWED_RESPONSE"
        return 1
    fi
    
    # Test device creation
    print_status "Testing device creation..."
    DEVICE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/devices \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"version_id\": \"$VERSION_ID\", \"name\": \"My Test Device\", \"dev_eui\": \"$DEV_EUI\", \"description\": \"Test device\"}")
    
    if echo "$DEVICE_RESPONSE" | grep -q "chirpstack_device_created"; then
        print_success "Device creation test passed"
    else
        print_error "Device creation test failed"
        echo "Response: $DEVICE_RESPONSE"
        return 1
    fi
    
    print_success "All API tests passed!"
}

# Main execution
main() {
    case "${1:-all}" in
        "unit")
            run_unit_tests
            ;;
        "integration")
            run_integration_tests
            ;;
        "api")
            run_api_tests
            ;;
        "all")
            print_status "Running all tests..."
            run_unit_tests && run_integration_tests && run_api_tests
            ;;
        *)
            echo "Usage: $0 [unit|integration|api|all]"
            echo "  unit        - Run unit tests only"
            echo "  integration - Run integration tests only"
            echo "  api         - Run API tests only"
            echo "  all         - Run all tests (default)"
            exit 1
            ;;
    esac
    
    if [ $? -eq 0 ]; then
        print_success "🎉 All tests completed successfully!"
    else
        print_error "❌ Some tests failed!"
        exit 1
    fi
}

# Run main function
main "$@"
