# 认证模块

本模块提供了完整的用户认证功能，包括用户注册、登录、JWT token生成和刷新等功能。

## 功能特性

- 用户注册和登录
- JWT Access Token 和 Refresh Token 支持
- 密码加密存储（使用bcrypt）
- 用户信息验证
- 自动token过期处理

## API 接口

### 1. 用户注册

**POST** `/api/v1/auth/register`

请求体：
```json
{
  "username": "newuser",
  "password": "password123",
  "email": "user@example.com",
  "nickname": "新用户"
}
```

响应：
```json
{
  "user": {
    "id": 1,
    "username": "newuser",
    "email": "user@example.com",
    "nickname": "新用户",
    "status": 1,
    "role_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. 用户登录

**POST** `/api/v1/auth/login`

请求体：
```json
{
  "username": "newuser",
  "password": "password123"
}
```

响应：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": 1,
    "username": "newuser",
    "email": "user@example.com",
    "nickname": "新用户",
    "status": 1,
    "role_id": 1
  }
}
```

### 3. 刷新Token

**POST** `/api/v1/auth/refresh`

请求体：
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

响应：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### 4. 获取用户信息

**GET** `/api/v1/auth/profile`

请求头：
```
Authorization: Bearer <access_token>
```

响应：
```json
{
  "id": 1,
  "username": "newuser",
  "email": "user@example.com",
  "nickname": "新用户",
  "status": 1,
  "role_id": 1,
  "last_login_at": "2024-01-01T00:00:00Z"
}
```

## 配置说明

在配置文件中设置JWT相关参数：

```yaml
jwt:
  secret: "your-secret-key-here"
  access_token_expire: 1    # access token过期时间(小时)
  refresh_token_expire: 168 # refresh token过期时间(小时，7天)
```

## 安全特性

1. **密码加密**：使用bcrypt算法加密存储密码
2. **Token分离**：Access Token和Refresh Token分离，提高安全性
3. **Token过期**：Access Token短期有效，Refresh Token长期有效
4. **输入验证**：严格的输入参数验证
5. **错误处理**：统一的错误响应格式

## 使用示例

### 使用curl进行测试

1. 注册用户：
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com",
    "nickname": "测试用户"
  }'
```

2. 用户登录：
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

3. 获取用户信息：
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <access_token>"
```

4. 刷新Token：
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

## 注意事项

1. 在生产环境中，请确保JWT Secret足够复杂且保密
2. 建议定期轮换JWT Secret
3. 可以根据业务需求调整Token过期时间
4. 密码强度验证可以根据需要调整
5. 建议在生产环境中使用HTTPS
