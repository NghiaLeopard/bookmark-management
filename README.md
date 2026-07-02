# bookmark-management

API quản lý bookmark — Go + Gin + Swagger.

## Cấu trúc project

```
bookmark-management/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── api/
│   │   └── api.go              # Gin engine, wiring routes
│   ├── config/
│   │   └── config.go           # Đọc biến môi trường (envconfig)
│   ├── handler/
│   │   ├── healthcheck.go
│   │   └── healtcheck_test.go  # Unit test handler
│   ├── service/
│   │   ├── healthcheck.go
│   │   ├── healthcheck_test.go # Unit test service
│   │   └── mocks/              # Mock generate bởi mockery
│   ├── model/
│   ├── initialize/
│   └── intergration_test/      # Integration test (full HTTP flow)
├── docs/                       # Swagger (swag init)
└── makefile
```

## Luồng chạy production

```
env (PORT, SERVICE_NAME, INSTANCE_ID)
  → config.NewConfig()
  → service.NewHealthCheck(cfg)
  → handler.NewHealthCheck(service)
  → Gin route
  → JSON response
```

## Chạy app

```bash
make run
```

Hoặc:

```bash
PORT=8080 SERVICE_NAME=bookmark-management INSTANCE_ID=1234567890 go run cmd/main.go
```

Swagger UI: `http://localhost:8080/swagger/index.html`

Generate Swagger:

```bash
swag init -g ./cmd/main.go -o ./docs
```

---

## Testing guide

### Nguyên tắc chung

**Không phải test nào cũng chạy y hệt luồng production.** Mỗi layer test một trách nhiệm, cô lập phần còn lại.

| Loại test | File | Test cái gì | Input | Có `t.Parallel`? |
|-----------|------|-------------|-------|------------------|
| Config unit | `config/config_test.go` | env → struct | `t.Setenv` | Không |
| Service unit | `service/*_test.go` | business logic | `&config.Config{...}` tay | Có |
| Handler unit | `handler/*_test.go` | HTTP status + JSON | mock service | Có |
| Integration | `intergration_test/*_test.go` | full stack wiring | `t.Setenv` + `api.NewEngine()` | Không |

### Test pyramid

```
        ┌─────────────────────┐
        │  Integration test   │  ít case — verify wire + HTTP end-to-end
        └─────────────────────┘
       ┌───────────────────────┐
       │   Handler unit test   │  mock service
       └───────────────────────┘
      ┌─────────────────────────┐
      │   Service unit test     │  struct config tay
      └─────────────────────────┘
     ┌───────────────────────────┐
     │   Config unit test        │  t.Setenv (tùy chọn)
     └───────────────────────────┘
```

---

### 1. Service unit test

**Mục tiêu:** test logic trong service (map config → response, nhánh if/else).

**Input:** tạo `&config.Config{...}` trực tiếp — **không dùng env**, **không gọi HTTP**.

```go
func TestCheckHealth(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name   string
        cfg    *config.Config
        wantSN string
        wantID string // rỗng = chỉ assert NotEmpty (UUID)
    }{
        {
            name:   "uses instance id from config",
            cfg:    &config.Config{ServiceName: "bookmark-management", InstanceId: "1234567890"},
            wantSN: "bookmark-management",
            wantID: "1234567890",
        },
        {
            name:   "generates uuid when instance id empty",
            cfg:    &config.Config{ServiceName: "bookmark-management", InstanceId: ""},
            wantSN: "bookmark-management",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            got := NewHealthCheck(tc.cfg).CheckHealth()

            assert.Equal(t, "OK", got.Message)
            assert.Equal(t, tc.wantSN, got.ServiceName)

            if tc.wantID != "" {
                assert.Equal(t, tc.wantID, got.InstanceId)
            } else {
                assert.NotEmpty(t, got.InstanceId)
            }
        })
    }
}
```

**Lưu ý:** không dùng biến global trong production code — mỗi lần gọi hàm phải có state riêng.

---

### 2. Handler unit test

**Mục tiêu:** test handler serialize HTTP đúng (status code + JSON body).

**Input:** mock service (mockery) — **không dùng config, không dùng env**.

```go
serviceMock := mocks.NewHealthCheck(t)
serviceMock.On("CheckHealth").Return(model.HealthCheck{
    Message: "OK", ServiceName: "bookmark-management", InstanceId: "1234567890",
})

handler := NewHealthCheck(serviceMock)
handler.CheckHealth(ctx)

assert.Equal(t, http.StatusOK, recorder.Code)
assert.Equal(t, `{"message":"OK","service_name":"bookmark-management","instance_id":"1234567890"}`, recorder.Body.String())
```

Generate mock:

```bash
go generate ./internal/service/...
```

---

### 3. Integration test

**Mục tiêu:** verify toàn bộ stack wire đúng — env → config → service → handler → HTTP response.

**Input:** `t.Setenv` + `api.NewEngine()` + `httptest`.

```go
func TestHealthCheckEP(t *testing.T) {
    // KHÔNG dùng t.Parallel() khi có t.Setenv

    testCases := []struct {
        name        string
        serviceName string
        instanceID  string
        expectUUID  bool
    }{
        {
            name:        "success",
            serviceName: "bookmark-management",
            instanceID:  "1234567890",
        },
        {
            name:        "generates uuid when instance id empty",
            serviceName: "bookmark-management",
            instanceID:  "",
            expectUUID:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Setenv("SERVICE_NAME", tc.serviceName)
            t.Setenv("INSTANCE_ID", tc.instanceID)

            recorder := httptest.NewRecorder()
            req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
            api.NewEngine().ServeHTTP(recorder, req)

            assert.Equal(t, http.StatusOK, recorder.Code)

            var got model.HealthCheck
            assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
            assert.Equal(t, "OK", got.Message)
            assert.Equal(t, tc.serviceName, got.ServiceName)

            if tc.expectUUID {
                assert.NotEmpty(t, got.InstanceId)
            } else {
                assert.Equal(t, tc.instanceID, got.InstanceId)
            }
        })
    }
}
```

**Quan trọng:**

- `api.NewEngine()` **tự gọi** `config.NewConfig()` — không cần tạo config tay trong integration test.
- **Không** dùng `t.Parallel()` cùng `t.Setenv` (Go sẽ panic).
- Assert so với **expected trong test case table**, không dùng `os.Getenv` để đối chiếu (case UUID random sẽ sai).
- `testing.T` chỉ có `Setenv` — đọc env trong code app dùng `os.Getenv`.

**Khi nào dùng `assert.JSONEq`:**

```go
assert.JSONEq(t, `{"message":"OK","service_name":"bookmark-management","instance_id":"1234567890"}`, recorder.Body.String())
```

Chỉ dùng khi mọi field deterministic (không có UUID random).

---

### 4. Config unit test (tùy chọn)

**Mục tiêu:** test riêng layer đọc env — tách khỏi service/integration.

```go
func TestNewConfig(t *testing.T) {
    t.Setenv("SERVICE_NAME", "bookmark-management")
    t.Setenv("INSTANCE_ID", "1234567890")

    cfg, err := config.NewConfig()
    assert.NoError(t, err)
    assert.Equal(t, "bookmark-management", cfg.ServiceName)
    assert.Equal(t, "1234567890", cfg.InstanceId)
}
```

Không cần assert config lại trong integration test nếu đã assert response cuối — response sai nghĩa là config hoặc layer giữa bị sai.

---

## Cách verify JSON response

| Layer | Cách verify |
|-------|-------------|
| Handler | So full JSON string (mock trả giá trị cố định) |
| Integration (deterministic) | `assert.JSONEq` hoặc unmarshal + assert từng field |
| Integration (UUID random) | Unmarshal + `assert.NotEmpty` cho field random |
| Service | Assert struct trực tiếp, không qua JSON |

---

## Chạy test

```bash
# Tất cả test
go test ./...

# Theo package
go test ./internal/service/... -v
go test ./internal/handler/... -v
go test ./internal/intergration_test/... -v

# Một test cụ thể
go test ./internal/service/... -run TestCheckHealth -v
```

---

## Checklist khi thêm feature mới

1. **Model** — struct + json tag trong `internal/model/`
2. **Service** — interface + implementation + `service/*_test.go` (config struct tay)
3. **Handler** — HTTP handler + swagger annotation + `handler/*_test.go` (mock service)
4. **API** — đăng ký route trong `internal/api/api.go`
5. **Integration** — `intergration_test/*_ep_test.go` (1–2 case chính, `t.Setenv`)
6. **Swagger** — chạy `swag init`, import `_ "…/docs"` trong `cmd/main.go`
7. **Mock** — `go generate` nếu thêm interface mới

---

## Quyết định nhanh: dùng gì?

```
Test logic nghiệp vụ?        → service test + config struct tay
Test HTTP serialize?         → handler test + mock
Test wire end-to-end?        → integration test + t.Setenv + NewEngine()
Test env đọc đúng?           → config test + t.Setenv
Cần chạy song song?          → service/handler: OK | integration/config với Setenv: KHÔNG
```
