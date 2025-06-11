# API Documentation - Go Authentication API

## Tổng quan

API này cung cấp hệ thống xác thực JWT hoàn chỉnh với các tính năng:
- Đăng ký người dùng
- Đăng nhập
- Xác thực JWT
- Bảo vệ routes
- Quản lý người dùng (CRUD)
- Tìm kiếm người dùng
- Phân trang

## Base URL
```
http://localhost:8080
```

## Authentication
API sử dụng JWT (JSON Web Token) để xác thực. Token có thời hạn 24 giờ.

### Header Format
```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### 1. Health Check

**GET** `/health`

Kiểm tra trạng thái API.

**Response:**
```json
{
  "status": "ok"
}
```

---

### 2. Đăng ký người dùng

**POST** `/api/v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Validation Rules:**
- `email`: Bắt buộc, định dạng email hợp lệ
- `password`: Bắt buộc, tối thiểu 6 ký tự
- `full_name`: Bắt buộc

**Success Response (201):**
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

**Error Response (400):**
```json
{
  "error": "user with email user@example.com already exists"
}
```

---

### 3. Đăng nhập

**POST** `/api/v1/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200):**
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

**Error Response (401):**
```json
{
  "error": "invalid credentials"
}
```

---

### 4. Profile (Protected)

**GET** `/api/v1/user/profile`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Success Response (200):**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "message": "Profile accessed successfully"
}
```

**Error Response (401):**
```json
{
  "error": "Authorization header required"
}
```

---

### 5. Lấy danh sách người dùng (Protected)

**GET** `/api/v1/users`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Query Parameters:**
- `page` (optional): Số trang (default: 1)
- `page_size` (optional): Số lượng user mỗi trang (default: 10, max: 100)

**Success Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "full_name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

---

### 6. Lấy thông tin người dùng theo ID (Protected)

**GET** `/api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Success Response (200):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "full_name": "John Doe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "User not found"
}
```

---

### 7. Cập nhật thông tin người dùng (Protected)

**PUT** `/api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "full_name": "New Full Name",
  "password": "newpassword123"
}
```

**Note:** Tất cả fields đều optional. Chỉ cần gửi fields muốn cập nhật.

**Success Response (200):**
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "full_name": "New Full Name",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

---

### 8. Xóa người dùng (Protected)

**DELETE** `/api/v1/users/:id`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Success Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

**Error Response (404):**
```json
{
  "error": "User not found"
}
```

---

### 9. Tìm kiếm người dùng (Protected)

**GET** `/api/v1/users/search`

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Query Parameters:**
- `q` (required): Từ khóa tìm kiếm (tìm trong email và full_name)
- `page` (optional): Số trang (default: 1)
- `page_size` (optional): Số lượng user mỗi trang (default: 10, max: 100)

**Success Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "full_name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 200  | Success |
| 201  | Created |
| 400  | Bad Request |
| 401  | Unauthorized |
| 500  | Internal Server Error |

## Testing Examples

### 1. Đăng ký người dùng mới
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### 2. Đăng nhập
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Truy cập Profile (cần token)
```bash
# Thay <TOKEN> bằng token thực từ response đăng nhập
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <TOKEN>"
```

### 4. Lấy danh sách người dùng
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer <TOKEN>"
```

### 5. Lấy thông tin người dùng theo ID
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <TOKEN>"
```

### 6. Cập nhật thông tin người dùng
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Name",
    "email": "updated@example.com"
  }'
```

### 7. Xóa người dùng
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <TOKEN>"
```

### 8. Tìm kiếm người dùng
```bash
curl -X GET "http://localhost:8080/api/v1/users/search?q=john&page=1&page_size=5" \
  -H "Authorization: Bearer <TOKEN>"
```

## Database Schema

### Users Table
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

1. **Password Hashing**: Sử dụng bcrypt với cost factor mặc định
2. **JWT Security**: Token có thời hạn 24 giờ
3. **Input Validation**: Validate tất cả input từ client
4. **SQL Injection Protection**: Sử dụng prepared statements
5. **CORS**: Có thể cấu hình CORS cho production

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=auth_db
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
```

## Production Considerations

1. **JWT Secret**: Thay đổi JWT_SECRET trong production
2. **Database**: Sử dụng managed database service
3. **HTTPS**: Luôn sử dụng HTTPS trong production
4. **Rate Limiting**: Implement rate limiting
5. **Logging**: Thêm structured logging
6. **Monitoring**: Setup health checks và monitoring
