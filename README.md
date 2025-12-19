# 012-go-api-auth-oauth2 — OAuth 2.0 & Single Sign-On (SSO)

โปรเจคนี้เป็นบทเรียน **OAuth 2.0 + Single Sign-On (SSO)** สำหรับ Go API  
โดยใช้ **Google OAuth 2.0** เป็นตัวอย่าง และออก **JWT ของระบบเราเอง** หลังจาก Login สำเร็จ

> แนวคิดนี้คือรูปแบบที่ใช้จริงในระบบ production  
> 👉 SSO Provider (Google) + Own JWT for Internal API

---

## 🎯 สิ่งที่ได้จากบทนี้

- OAuth 2.0 Authorization Code Flow
- Login with Google
- Google Callback + Exchange Token
- Fetch Google UserInfo (OpenID Connect)
- Upsert User + OAuth Identity ใน PostgreSQL
- Generate JWT (Access Token) ของระบบเราเอง
- Protected API ด้วย JWT Middleware
- โครงสร้างโปรเจคแบบ Production-grade

---

## 🧱 Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- OAuth 2.0 (Google)
- JWT (HS256)
- Docker / Docker Compose

---

## 📂 Project Structure

```
012-go-api-auth-oauth2/
├─ cmd/
│  └─ api/
│     └─ main.go
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ db/
│  │  └─ postgres.go
│  ├─ models/
│  │  ├─ user.go
│  │  └─ oauth_identity.go
│  ├─ repositories/
│  │  ├─ user_repo.go
│  │  └─ identity_repo.go
│  ├─ services/
│  │  ├─ jwt_service.go
│  │  └─ oauth_google_service.go
│  ├─ middlewares/
│  │  └─ auth_middleware.go
│  ├─ handlers/
│  │  └─ oauth_handler.go
│  ├─ routes/
│  │  └─ routes.go
│  └─ utils/
│     └─ response.go
├─ .env.example
├─ docker-compose.yml
├─ go.mod
└─ README.md
```

---

## 🌐 API Endpoints

### Public

- GET /auth/google/login
- GET /auth/google/callback

### Protected

- GET /me

---

## ⚙️ Environment Configuration

```env
APP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_oauth2

JWT_ISSUER=012-go-api-auth-oauth2
JWT_ACCESS_SECRET=CHANGE_ME_ACCESS

GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET=YOUR_GOOGLE_CLIENT_SECRET
GOOGLE_REDIRECT_URL=http://localhost:8081/auth/google/callback

FRONTEND_SUCCESS_REDIRECT=http://localhost:3000/oauth/success
```

---

## 🐳 Docker Compose

```yaml
services:
  postgres:
    image: postgres:16
    ports:
      - "5432:5432"
```

---

## 🚀 Getting Started

```bash
docker compose up -d
go mod init github.com/ChinawatDc/012-go-api-auth-oauth2
go mod tidy
go run ./cmd/api
```

---

## 🔐 OAuth Flow

1. Client → /auth/google/login
2. Redirect to Google
3. Google → /auth/google/callback
4. Exchange code → token
5. Fetch userinfo
6. Upsert user + identity
7. Generate JWT
8. Return / Redirect

---

## 🛡 Security

- Authorization Code Flow
- OAuth State (CSRF Protection)
- Own JWT for internal APIs
- No Google token storage

---

## 📚 Next Lesson

**013-go-api-auth-ldap** — LDAP / Active Directory Authentication

---

Author: Chinawat Daochai  
Course: Mastering Go API Development
