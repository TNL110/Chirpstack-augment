# External API Access

This document describes how to access the Go Auth API from external servers and networks.

## Configuration

The API is configured to accept connections from external networks by binding to all network interfaces (`0.0.0.0`).

### Docker Compose Configuration

```yaml
services:
  api:
    build: .
    ports:
      - "0.0.0.0:8080:8080"  # Bind to all interfaces
    # ... other configuration

  postgres:
    image: postgres:15-alpine
    ports:
      - "0.0.0.0:5432:5432"  # PostgreSQL also accessible externally
    # ... other configuration
```

## Server Information

- **Server IP**: `192.168.0.93`
- **API Port**: `8080`
- **Database Port**: `5432` (if needed)

## API Endpoints

### Base URL
```
http://192.168.0.93:8080
```

### Available Endpoints

#### Public Endpoints (No Authentication Required)

1. **Health Check**
   ```bash
   GET http://192.168.0.93:8080/health
   ```
   Response:
   ```json
   {"status":"ok"}
   ```

2. **User Registration**
   ```bash
   POST http://192.168.0.93:8080/api/v1/auth/register
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "password123",
     "full_name": "User Name"
   }
   ```

3. **User Login**
   ```bash
   POST http://192.168.0.93:8080/api/v1/auth/login
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "password123"
   }
   ```

#### Protected Endpoints (Require JWT Token)

1. **Get User Profile**
   ```bash
   GET http://192.168.0.93:8080/api/v1/user/profile
   Authorization: Bearer <jwt_token>
   ```

2. **Get All Users**
   ```bash
   GET http://192.168.0.93:8080/api/v1/users?page=1&page_size=10
   Authorization: Bearer <jwt_token>
   ```

3. **Get User by ID**
   ```bash
   GET http://192.168.0.93:8080/api/v1/users/{id}
   Authorization: Bearer <jwt_token>
   ```

4. **Update User**
   ```bash
   PUT http://192.168.0.93:8080/api/v1/users/{id}
   Authorization: Bearer <jwt_token>
   Content-Type: application/json
   
   {
     "email": "newemail@example.com",
     "full_name": "New Name"
   }
   ```

5. **Delete User**
   ```bash
   DELETE http://192.168.0.93:8080/api/v1/users/{id}
   Authorization: Bearer <jwt_token>
   ```

6. **Search Users**
   ```bash
   GET http://192.168.0.93:8080/api/v1/users/search?q=search_term&page=1&page_size=10
   Authorization: Bearer <jwt_token>
   ```

## Testing External Access

### Automated Testing

Use the provided test script to verify external access:

```bash
./test_external_api.sh
```

This script tests:
- Health check
- User registration (with ChirpStack integration)
- User login
- Protected routes with JWT authentication
- Unauthorized access handling

### Manual Testing Examples

#### 1. Register a New User
```bash
curl -X POST http://192.168.0.93:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

#### 2. Login
```bash
curl -X POST http://192.168.0.93:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 3. Access Protected Route
```bash
# Use token from login response
curl -X GET http://192.168.0.93:8080/api/v1/user/profile \
  -H "Authorization: Bearer <your_jwt_token>"
```

## ChirpStack Integration

When registering users externally, the ChirpStack integration automatically creates:

- **Tenant** with user's email as name
- **Application** named "Lnode"
- **Device Profile** named "RAK_ABP" with custom payload decoder

This happens automatically for all user registrations when ChirpStack is enabled.

## Security Considerations

### Authentication
- All protected endpoints require valid JWT tokens
- Tokens are obtained through login or registration
- Tokens have expiration time (24 hours by default)

### Network Security
- API is accessible from any IP address
- Consider implementing IP whitelisting if needed
- Use HTTPS in production environments
- Ensure firewall rules are properly configured

### Database Access
- PostgreSQL is also exposed externally on port 5432
- Use strong passwords and consider restricting database access
- Monitor database connections

## Troubleshooting

### Connection Issues

1. **Cannot connect to API**
   - Check if server is running: `docker ps`
   - Verify port is open: `netstat -tlnp | grep :8080`
   - Check firewall settings
   - Verify server IP address

2. **Authentication Errors**
   - Ensure JWT token is valid and not expired
   - Check Authorization header format: `Bearer <token>`
   - Verify token was obtained from login/register

3. **Network Issues**
   - Test connectivity: `ping 192.168.0.93`
   - Check if port is accessible: `telnet 192.168.0.93 8080`
   - Verify no proxy/firewall blocking

### Logs

Check application logs:
```bash
docker logs go-auth-api --tail=50
```

Check all services:
```bash
docker-compose logs
```

## Production Deployment

For production deployment, consider:

1. **HTTPS/TLS**: Use reverse proxy (nginx) with SSL certificates
2. **Domain Name**: Use proper domain instead of IP address
3. **Load Balancing**: Multiple API instances behind load balancer
4. **Monitoring**: Implement health checks and monitoring
5. **Security**: IP whitelisting, rate limiting, WAF
6. **Backup**: Regular database backups
7. **Environment Variables**: Secure secret management
