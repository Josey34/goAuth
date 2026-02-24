# 🔐 GoAuth — Authentication Layer (Project 2 of 3)

> Add authentication, JWT tokens, and role-based access control to DocVault — protecting document routes so only registered users can upload, download, and manage files.

## 🎯 System Overview

This is **Project 2** of a 3-project continuous system:

```
┌──────────────────────────────────────────────────────────────────┐
│                  Document Management System                      │
│                                                                  │
│  Project 1: DocVault ✅      ──► Core API (done!)                │
│                                                                  │
│  Project 2: GoAuth           ──► Authentication layer             │
│  (You are here)                  • User registration & login     │
│                                  • JWT tokens (access + refresh) │
│                                  • Protect DocVault routes       │
│                                  • Users own their documents     │
│                                  • Admin role (manage users)     │
│                                                                  │
│  Project 3: GoFlow           ──► Document processor               │
│                                                                  │
│  Shared: SQLite DB (users + documents tables)                    │
└──────────────────────────────────────────────────────────────────┘
```

### What changes in DocVault after GoAuth:
- All document routes require `Authorization: Bearer <token>` header
- Documents get an `uploaded_by` field linking to the user
- Users can only delete their own files (admins can delete any)
- New admin routes for user management

## 🧰 Tech Stack

| Tool | Purpose | Why This? |
|------|---------|-----------|
| **Go (Gin)** | HTTP framework | Same as DocVault |
| **SQLite** | Database | Same DB file — shares `users` table with DocVault |
| **golang-jwt/jwt** | JWT tokens | Industry standard |
| **crypto/bcrypt** | Password hashing | Go standard library |
| **testify** | Testing | Mock interfaces, assertions |

---

## 📚 Go Concepts You Will Learn

### Authentication & Security (NEW)
- [ ] Password hashing with `crypto/bcrypt` (cost factor, salt)
- [ ] JWT structure (header, payload, signature)
- [ ] JWT signing with HMAC (HS256)
- [ ] JWT claims (sub, role, exp, iat, type)
- [ ] Token validation and expiry checking
- [ ] Access tokens (short-lived) vs refresh tokens (long-lived)
- [ ] Token rotation on refresh
- [ ] Auth middleware (extract + validate token)
- [ ] Role-based access control (admin vs user)
- [ ] Preventing email enumeration (same error for wrong email & wrong password)

### Custom Error Types (NEW)
- [ ] Custom error structs implementing `error` interface
- [ ] Error domain design (AuthError, ValidationError, NotFoundError, ConflictError)
- [ ] `errors.As()` type assertions
- [ ] Centralized error → HTTP status code mapping

### Unit Testing (NEW depth)
- [ ] `testing` package, `t.Run()` subtests
- [ ] Table-driven tests (Go's idiomatic pattern)
- [ ] Mock interfaces with testify/mock
- [ ] Testing usecases with mock repo + mock token service
- [ ] Testing handlers with `httptest`
- [ ] Testing middleware (valid/invalid/expired tokens, role checks)
- [ ] `go test -cover`, `go test -race`

### Clean Architecture (Reinforced + Deeper)
- [ ] Entity layer (User — NO password hash in responses)
- [ ] Repository pattern + service interfaces
- [ ] Usecase layer (AuthUsecase, UserUsecase)
- [ ] Factory pattern with dependency injection
- [ ] **Testability through interfaces** — mock any layer

### Middleware Chains (NEW)
- [ ] Auth middleware (validate JWT, inject user into Gin context)
- [ ] Role middleware (check role, reject unauthorized)
- [ ] `c.Set()` / `c.Get()` for passing data through chain
- [ ] `c.Abort()` to stop chain on failure
- [ ] Nested route groups (public → protected → admin)

---

## 📁 Project Structure

```
goauth/
├── main.go                     # Entry point: factory → server (integrates with DocVault)
├── go.mod / go.sum
├── .env                        # JWT_SECRET, DB_PATH (same DB as DocVault!), TTLs
│
├── config/
│   └── config.go               # JWT secret, access/refresh TTLs, bcrypt cost, DB path, port
│
├── factory/
│   └── factory.go              # Creates all deps: repo, token service, usecases, handlers, middleware
│
├── entity/
│   └── user.go                 # Domain struct — NO JSON tags: ID, Email, Name, PasswordHash, Role, CreatedAt
│
├── dto/
│   ├── request.go              # RegisterRequest, LoginRequest, UpdateProfileRequest, RefreshRequest
│   └── response.go             # UserResponse (NO password_hash!), TokenResponse, ErrorResponse, FromEntity()
│
├── repository/
│   ├── repository.go           # Interface: UserRepository (Insert, FindByEmail, FindByID, FindAll, Delete, Update)
│   └── sqlite_user.go          # SQLite impl — uses the SAME DB file as DocVault
│
├── service/
│   ├── token.go                # Interface: TokenService (GenerateAccess, GenerateRefresh, Validate)
│   └── token_jwt.go            # JWT implementation
│
├── usecase/
│   ├── auth.go                 # AuthUsecase: Register, Login, RefreshToken
│   └── user.go                 # UserUsecase: GetProfile, UpdateProfile, ListUsers, DeleteUser
│
├── handler/
│   ├── auth.go                 # POST /register, /login, /refresh
│   └── user.go                 # GET /me, PUT /me, GET /admin/users, DELETE /admin/users/:id
│
├── middleware/
│   ├── auth.go                 # JWT validation → inject user into context
│   └── role.go                 # Role check → allow/deny
│
├── errors/
│   └── errors.go               # AuthError, ValidationError, NotFoundError, ConflictError + ToHTTPStatus()
│
├── database/
│   └── sqlite.go               # Opens SAME SQLite file as DocVault
│
└── tests/
    ├── mock/
    │   ├── mock_user_repo.go
    │   └── mock_token_svc.go
    ├── usecase/
    │   ├── auth_test.go        # Register, Login, Refresh — happy + error paths
    │   └── user_test.go        # Profile, ListUsers, Delete — happy + error paths
    ├── handler/
    │   ├── auth_test.go
    │   └── user_test.go
    ├── middleware/
    │   ├── auth_test.go        # Valid/invalid/expired/missing tokens
    │   └── role_test.go        # Admin allowed, user denied
    └── service/
        └── token_test.go       # Generate, validate, expiry, type checking
```

### How the layers connect

```
HTTP Request
    │
    ▼
middleware/auth.go       → Extract JWT → validate via TokenService → inject user into context
    │
    ▼
middleware/role.go       → Check role from context → allow/deny
    │
    ▼
handler/                 → Parse DTO → call usecase → return DTO response
    │
    ▼
usecase/                 → Business logic (hash, verify, orchestrate)
    │
    ├──► repository.UserRepository  (interface) → sqlite_user.go → SQLite
    └──► service.TokenService       (interface) → token_jwt.go → JWT

Factory wires everything. main.go stays clean.
```

### Factory Pattern

```go
type Factory struct {
    AuthHandler    *handler.AuthHandler
    UserHandler    *handler.UserHandler
    AuthMiddleware *middleware.AuthMiddleware
    RoleMiddleware *middleware.RoleMiddleware
}

func New(cfg *config.Config) (*Factory, error) {
    db, _ := database.OpenSQLite(cfg.DBPath) // SAME DB as DocVault!
    userRepo := repository.NewSQLiteUserRepo(db)
    tokenSvc := service.NewJWTTokenService(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
    authUC := usecase.NewAuthUsecase(userRepo, tokenSvc, cfg.BcryptCost)
    userUC := usecase.NewUserUsecase(userRepo)
    return &Factory{
        AuthHandler:    handler.NewAuthHandler(authUC),
        UserHandler:    handler.NewUserHandler(userUC),
        AuthMiddleware: middleware.NewAuthMiddleware(tokenSvc),
        RoleMiddleware: middleware.NewRoleMiddleware(),
    }, nil
}
```

---

## 🔌 API Endpoints

### Public Routes (no auth required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/register` | Create account (email, password, name) → user + tokens |
| `POST` | `/api/auth/login` | Login → tokens |
| `POST` | `/api/auth/refresh` | Refresh expired access token |

### Protected Routes (valid JWT required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users/me` | Get own profile |
| `PUT` | `/api/users/me` | Update own profile |

### Admin Routes (admin role required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/admin/users` | List all users |
| `DELETE` | `/api/admin/users/:id` | Delete a user |

### DocVault Routes (NOW PROTECTED after GoAuth integration)
| Method | Endpoint | Auth | Change |
|--------|----------|------|--------|
| `POST` | `/api/documents/upload` | Required | Added `uploaded_by` = userID from token |
| `GET` | `/api/documents` | Required | Filter by own docs (admin sees all) |
| `DELETE` | `/api/documents/:id` | Required | Only own docs (admin can delete any) |
| `GET` | `/api/documents/:id/download` | Required | Only own docs (admin can download any) |

---

## 🗺️ Phase-by-Phase Roadmap

### Phase 1: Project Setup — Shared DB (Day 1)

**Goal:** Set up GoAuth pointing to the SAME SQLite DB as DocVault.

**Steps:**
1. `go mod init goauth`
2. Create `.env` — `DB_PATH` points to DocVault's SQLite file
3. `config/config.go` — JWT secret, TTLs, bcrypt cost, DB path
4. `database/sqlite.go` — opens same DB (users table already exists from DocVault's migrations)
5. `entity/user.go` — pure domain struct (ZERO imports)
6. `errors/errors.go` — AuthError, ValidationError, NotFoundError, ConflictError + `ToHTTPStatus()`
7. `factory/factory.go` (skeleton), `main.go`
8. `GET /ping` → pong

**Test yourself:**
- [ ] GoAuth connects to DocVault's SQLite file
- [ ] `entity/user.go` has ZERO external imports

---

### Phase 2: Registration (Day 2)

**Goal:** User registration with bcrypt, following clean architecture.

**Steps:**
1. `repository/repository.go` — `UserRepository` interface + `sqlite_user.go`
2. `usecase/auth.go` — `AuthUsecase.Register()`: validate → check email exists → hash → save
3. `dto/request.go` + `dto/response.go` — `UserResponse` has NO password_hash
4. `handler/auth.go` — Register handler
5. Wire in factory, register route

**Test yourself:**
- [ ] Register works, password is bcrypt hash in DB
- [ ] Duplicate email → 409
- [ ] Response has NO password_hash
- [ ] Usecase does NOT import `database/sql`

---

### Phase 3: Login & JWT (Day 3)

**Goal:** Verify password, issue JWT tokens via token service interface.

**Steps:**
1. `service/token.go` — `TokenService` interface (GenerateAccess, GenerateRefresh, Validate)
2. `service/token_jwt.go` — JWT implementation with claims: sub, role, exp, iat, type
3. Update `AuthUsecase.Login()` — find user → bcrypt compare → generate tokens
4. **SECURITY:** Same error for wrong email AND wrong password: `"invalid email or password"`

**Test yourself:**
- [ ] Login → access_token + refresh_token
- [ ] Wrong email AND wrong password → identical 401 error
- [ ] Decode at jwt.io → claims correct
- [ ] Usecase uses `TokenService` interface, NOT `golang-jwt`

---

### Phase 4: Auth Middleware (Day 4)

**Goal:** Middleware validates JWT and injects user into Gin context.

**Steps:**
1. `middleware/auth.go` — depends on `TokenService` interface
   - Extract `Authorization: Bearer <token>` → validate → `c.Set("userID", id)` / `c.Set("userRole", role)` → or 401 + `c.Abort()`
2. `usecase/user.go` — `UserUsecase.GetProfile(ctx, userID)`
3. `handler/user.go` — GetProfile handler (reads userID from context)
4. Create protected route group in main.go

**Test yourself:**
- [ ] No token → 401
- [ ] Invalid token → 401
- [ ] Valid token → 200 with profile

---

### Phase 5: Profile Management (Day 5)

**Goal:** View and update own profile.

**Steps:** Add Update to repo → UpdateProfile to usecase → UpdateProfile handler

**Test yourself:**
- [ ] Update name works
- [ ] Can only update own profile

---

### Phase 6: Role-Based Access Control (Day 6)

**Goal:** Admin-only routes.

**Steps:**
1. `middleware/role.go` — `RequireRole("admin")` → 403 + `c.Abort()` if wrong role
2. ListUsers + DeleteUser usecase methods + handlers
3. Admin route group with both auth + role middleware
4. Admin cannot delete themselves

**Test yourself:**
- [ ] Regular user → admin routes → 403
- [ ] Admin → admin routes → 200
- [ ] Admin cannot delete self

---

### Phase 7: Token Refresh (Day 7)

**Goal:** Refresh expired access tokens.

**Steps:** `AuthUsecase.RefreshToken()` — validate → check type is "refresh" → verify user exists → generate new pair (token rotation)

**Test yourself:**
- [ ] Refresh → new token pair
- [ ] Access token as refresh → fails
- [ ] Deleted user refresh → fails

---

### Phase 8: Integrate Auth into DocVault (Day 8)

**Goal:** Protect DocVault's document routes with GoAuth's middleware.

**Steps:**
1. **Add `uploaded_by` column** to documents table (migration)
2. **Copy auth + role middleware** into DocVault (or import as shared package)
3. **Update DocVault's factory** — add TokenService and middleware
4. **Protect routes:**
   ```go
   api := r.Group("/api")
   api.Use(authMiddleware.Required())
   {
       api.POST("/documents/upload", docHandler.Upload)  // sets uploaded_by from token
       api.GET("/documents", docHandler.List)             // filters by user (admin sees all)
       api.DELETE("/documents/:id", docHandler.Delete)    // own docs only (admin: any)
   }
   ```
5. **Update document usecase** — accept userID, check ownership on delete/download
6. **Update document entity** — add `UploadedBy` field

**Test yourself:**
- [ ] Upload without token → 401
- [ ] Upload with token → document has `uploaded_by` = your userID
- [ ] Delete someone else's doc → 403 (admin can delete any)
- [ ] List shows only your docs (admin sees all)

---

### Phase 9: Unit Tests (Day 9–10)

**Goal:** Comprehensive tests for all layers.

**Steps:**
1. **Mocks:** `MockUserRepo`, `MockTokenService`
2. **Usecase tests** (table-driven):
   - Register: success, duplicate email, invalid input
   - Login: success, wrong password, user not found (SAME error!)
   - Refresh: success, wrong token type, deleted user, expired
3. **Handler tests** with `httptest`
4. **Middleware tests:** valid/invalid/expired tokens, admin/user roles
5. **Token service tests:** generate, validate, expiry, type check

```go
func TestLogin(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        password string
        setup   func(*mock.MockUserRepo)
        wantErr bool
    }{
        {"success", "test@test.com", "pass123", func(m *mock.MockUserRepo) {
            m.On("FindByEmail", mock.Anything, "test@test.com").Return(&entity.User{PasswordHash: hashedPass}, nil)
        }, false},
        {"wrong password", "test@test.com", "wrong", func(m *mock.MockUserRepo) {
            m.On("FindByEmail", mock.Anything, "test@test.com").Return(&entity.User{PasswordHash: hashedPass}, nil)
        }, true},
        {"user not found", "no@test.com", "pass", func(m *mock.MockUserRepo) {
            m.On("FindByEmail", mock.Anything, "no@test.com").Return(nil, errors.NotFoundError{})
        }, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) { ... })
    }
}
```

**Test yourself:**
- [ ] `go test ./...` passes
- [ ] `go test -race ./...` clean
- [ ] Each error case has its own test
- [ ] ZERO real DB calls in usecase tests

---

### Phase 10: Polish (Day 11)

**Goal:** Error handling, validation, input sanitization.

**Steps:**
1. Centralized `errors.ToHTTPStatus()` mapping
2. Input validation in usecase: email format, password length, trim whitespace, lowercase email
3. Consistent error format everywhere
4. Request logging middleware

**Test yourself:**
- [ ] `" Josey@EXAMPLE.com "` stored as `"josey@example.com"`
- [ ] All errors: `{"error": "message"}`

---

## 🧪 Testing Cheat Sheet

```bash
go run main.go

# Register
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"josey@test.com","password":"securepass123","name":"Josey"}'

# Login
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"josey@test.com","password":"securepass123"}'

TOKEN="eyJ..."
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/users/me
curl -X POST http://localhost:8081/api/auth/refresh -H "Content-Type: application/json" -d '{"refresh_token":"eyJ..."}'

# After DocVault integration:
curl -H "Authorization: Bearer $TOKEN" -F "file=@test.pdf" http://localhost:8080/api/documents/upload

# Tests
go test ./... && go test -race ./... && go test -cover ./...
```

## 📖 Dependencies
```bash
go get github.com/gin-gonic/gin github.com/mattn/go-sqlite3 github.com/golang-jwt/jwt/v5 github.com/google/uuid github.com/stretchr/testify
```

## 💡 Tutor Instructions
1. Don't give full solutions — function signatures only
2. Ask security questions: "why same error for wrong email and wrong password?"
3. Point me to jwt.io for debugging tokens
4. Check architecture: handler must not import repository
5. Guide testing: "what edge cases should this cover?"

### Common mistakes:
- Exposing password_hash in responses (use DTO!)
- Different errors for wrong email vs password (email enumeration risk)
- Using access token as refresh token
- JWT secret in code instead of env
- Not lowercasing email
- Forgetting `c.Abort()` in middleware
- JSON tags on entity structs
- Not testing error paths

---

## ✅ Completion Checklist

- [ ] Phase 1: Setup, shared DB, entity, errors, factory skeleton
- [ ] Phase 2: Registration with bcrypt
- [ ] Phase 3: Login with JWT via TokenService interface
- [ ] Phase 4: Auth middleware
- [ ] Phase 5: Profile management
- [ ] Phase 6: Role-based access control
- [ ] Phase 7: Token refresh
- [ ] Phase 8: **Integrate auth into DocVault** (protect routes, document ownership)
- [ ] Phase 9: Unit tests (usecase, handler, middleware, token service)
- [ ] Phase 10: Error mapping, validation, polish

**→ Next: Project 3 (GoFlow) — Process documents uploaded through DocVault**