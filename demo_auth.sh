#!/bin/bash

# 认证API演示脚本
# 确保服务器正在运行在 localhost:8080

BASE_URL="http://localhost:8080/api/v1/auth"

echo "=== 工单系统认证API演示 ==="
echo

# 1. 用户注册
echo "1. 用户注册"
echo "POST $BASE_URL/register"
curl -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo_user",
    "password": "password123",
    "email": "demo@example.com",
    "nickname": "演示用户"
  }' | jq .
echo
echo

# 2. 用户登录
echo "2. 用户登录"
echo "POST $BASE_URL/login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo_user",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .

# 提取access token
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token')

echo
echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"
echo

# 3. 获取用户信息
echo "3. 获取用户信息"
echo "GET $BASE_URL/profile"
curl -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo
echo

# 4. 刷新Token
echo "4. 刷新Token"
echo "POST $BASE_URL/refresh"
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "$REFRESH_RESPONSE" | jq .

# 提取新的access token
NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.access_token')
NEW_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.refresh_token')

echo
echo "New Access Token: $NEW_ACCESS_TOKEN"
echo "New Refresh Token: $NEW_REFRESH_TOKEN"
echo

# 5. 使用新token获取用户信息
echo "5. 使用新token获取用户信息"
echo "GET $BASE_URL/profile"
curl -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN" | jq .
echo
echo

echo "=== 演示完成 ==="
