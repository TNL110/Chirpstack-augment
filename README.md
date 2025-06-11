# Go Authentication API

Một REST API được xây dựng bằng Go với hệ thống xác thực JWT, PostgreSQL và Docker.

## Tính năng

- ✅ Đăng ký người dùng (Register)
- ✅ Đăng nhập (Login)
- ✅ Xác thực JWT
- ✅ Bảo vệ routes với middleware
- ✅ Quản lý người dùng (CRUD operations)
- ✅ Tìm kiếm người dùng
- ✅ Phân trang (Pagination)
- ✅ PostgreSQL database
- ✅ Docker deployment
- ✅ Password hashing với bcrypt
- ✅ **ChirpStack Integration** - Tự động tạo ChirpStack resources khi đăng ký user

## Cấu trúc Project

```
go-auth-api/
├── cmd/server/          # Entry point
├── internal/
│   ├── auth/           # JWT và password utilities
│   ├── config/         # Configuration
│   ├── database/       # Database connection
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Middleware functions
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic
├── docker-compose.yml  # Docker services
├── Dockerfile         # API container
├── init.sql          # Database schema
└── .env              # Environment variables
```

## Cài đặt và Chạy

### 1. Sử dụng Docker (Khuyến nghị)

```bash
# Clone repository
git clone <repository-url>
cd go-auth-api

# Chạy với Docker Compose
docker-compose up --build

# API sẽ chạy tại: http://localhost:8080
```

### 2. Chạy Local (cần Go và PostgreSQL)

```bash
# Cài đặt dependencies
go mod tidy

# Chạy PostgreSQL local hoặc update .env file

# Chạy application
go run cmd/server/main.go
```

## API Endpoints

### Base URL: `http://localhost:8080/api/v1`

### 1. Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "ok"
}
```

### 2. Đăng ký (Register)
```
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. Đăng nhập (Login)
```
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4. Profile (Protected)
```
GET /api/v1/user/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "message": "Profile accessed successfully"
}
```

### 5. Quản lý người dùng (Protected)

#### Lấy danh sách người dùng
```
GET /api/v1/users?page=1&page_size=10
Authorization: Bearer <token>
```

#### Lấy thông tin người dùng theo ID
```
GET /api/v1/users/:id
Authorization: Bearer <token>
```

#### Cập nhật thông tin người dùng
```
PUT /api/v1/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "full_name": "New Name",
  "password": "newpassword123"
}
```

#### Xóa người dùng
```
DELETE /api/v1/users/:id
Authorization: Bearer <token>
```

#### Tìm kiếm người dùng
```
GET /api/v1/users/search?q=keyword&page=1&page_size=10
Authorization: Bearer <token>
```

## Testing với cURL

### 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```

### 2. Đăng ký
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 3. Đăng nhập
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Truy cập Profile (thay <TOKEN> bằng token thực)
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <TOKEN>"
```

## ChirpStack Integration

API hỗ trợ tích hợp với ChirpStack để tự động tạo resources khi user đăng ký.

### Cấu hình ChirpStack

Thêm các environment variables sau:

```env
CHIRPSTACK_ENABLED=true
CHIRPSTACK_HOST=192.168.0.21
CHIRPSTACK_PORT=8090
CHIRPSTACK_TOKEN=your-chirpstack-api-token
```

### Tính năng

Khi user đăng ký, hệ thống tự động tạo:
- **Tenant** với tên là email của user
- **Application** tên "Lnode" trong tenant
- **Device Profile** tên "RAK_ABP" với payload decoder

### Test ChirpStack Integration

```bash
./test_chirpstack_integration.sh
```

Chi tiết xem: [CHIRPSTACK_INTEGRATION.md](CHIRPSTACK_INTEGRATION.md)

## External Access

API có thể được truy cập từ các server khác qua mạng:

- **Server IP**: `192.168.0.93`
- **API URL**: `http://192.168.0.93:8080`
- **Health Check**: `http://192.168.0.93:8080/health`

### Test External Access

```bash
./test_external_api.sh
```

Chi tiết xem: [EXTERNAL_ACCESS.md](EXTERNAL_ACCESS.md)

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=auth_db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080

# ChirpStack Configuration (Optional)
CHIRPSTACK_ENABLED=true
CHIRPSTACK_HOST=192.168.0.21
CHIRPSTACK_PORT=8090
CHIRPSTACK_TOKEN=your-chirpstack-api-token
```

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Security Features

- ✅ Password hashing với bcrypt
- ✅ JWT tokens với expiration
- ✅ Protected routes với middleware
- ✅ Input validation
- ✅ SQL injection protection với prepared statements

## Troubleshooting

### Lỗi Database Connection
```bash
# Kiểm tra PostgreSQL container
docker-compose logs postgres

# Restart services
docker-compose down
docker-compose up --build
```

### Lỗi Port đã được sử dụng
```bash
# Tìm process sử dụng port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```
