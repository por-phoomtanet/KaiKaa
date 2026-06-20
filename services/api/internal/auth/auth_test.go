package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// newTestDB ต่อ mongo ผ่าน MONGO_URI (default localhost) ใช้ db ชั่วคราวแล้ว drop ทิ้ง
// ถ้าต่อ mongo ไม่ได้ จะ t.Skip — CI ที่ไม่มี mongo จะไม่ fail
func newTestDB(t *testing.T) *mongo.Database {
	t.Helper()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("mongo unavailable (%s): %v", uri, err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("mongo ping failed (%s): %v", uri, err)
	}

	db := client.Database("kaikaa_test_" + primitive.NewObjectID().Hex())
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Drop(ctx)
		_ = client.Disconnect(ctx)
	})
	return db
}

func newRouter(db *mongo.Database) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewHandler(NewService(db, "test-secret"))
	r := gin.New()
	g := r.Group("/api/auth")
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	return r
}

func do(r *gin.Engine, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type envelope struct {
	Success bool `json:"success"`
	Data    struct {
		Token string `json:"token"`
		Shop  struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"shop"`
	} `json:"data"`
}

func TestRegister(t *testing.T) {
	r := newRouter(newTestDB(t))
	body := `{"shop_name":"ร้านกาแฟ","email":"a@b.com","password":"secret123"}`

	w := do(r, "/api/auth/register", body)
	if w.Code != http.StatusOK {
		t.Fatalf("register code = %d, want 200 (body=%s)", w.Code, w.Body)
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if !env.Success || env.Data.Token == "" {
		t.Errorf("expected success + token, got %+v", env)
	}
	if env.Data.Shop.Name != "ร้านกาแฟ" {
		t.Errorf("shop name = %q, want UTF-8 Thai preserved", env.Data.Shop.Name)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	r := newRouter(newTestDB(t))
	body := `{"shop_name":"X","email":"dup@b.com","password":"secret123"}`

	if w := do(r, "/api/auth/register", body); w.Code != http.StatusOK {
		t.Fatalf("first register code = %d, want 200", w.Code)
	}
	if w := do(r, "/api/auth/register", body); w.Code != http.StatusConflict {
		t.Errorf("duplicate register code = %d, want 409", w.Code)
	}
}

func TestRegisterValidation(t *testing.T) {
	r := newRouter(newTestDB(t))
	cases := map[string]string{
		"missing shop_name": `{"email":"a@b.com","password":"secret123"}`,
		"invalid email":     `{"shop_name":"X","email":"not-email","password":"secret123"}`,
		"short password":    `{"shop_name":"X","email":"a@b.com","password":"123"}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			if w := do(r, "/api/auth/register", body); w.Code != http.StatusBadRequest {
				t.Errorf("code = %d, want 400", w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	r := newRouter(newTestDB(t))
	reg := `{"shop_name":"Cafe","email":"login@b.com","password":"secret123"}`
	if w := do(r, "/api/auth/register", reg); w.Code != http.StatusOK {
		t.Fatalf("setup register failed: %d", w.Code)
	}

	t.Run("correct password", func(t *testing.T) {
		w := do(r, "/api/auth/login", `{"email":"login@b.com","password":"secret123"}`)
		if w.Code != http.StatusOK {
			t.Fatalf("code = %d, want 200 (body=%s)", w.Code, w.Body)
		}
		var env envelope
		_ = json.Unmarshal(w.Body.Bytes(), &env)
		if env.Data.Token == "" {
			t.Error("expected token on successful login")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		w := do(r, "/api/auth/login", `{"email":"login@b.com","password":"WRONG"}`)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code = %d, want 401", w.Code)
		}
	})

	t.Run("unknown email", func(t *testing.T) {
		w := do(r, "/api/auth/login", `{"email":"ghost@b.com","password":"secret123"}`)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("code = %d, want 401", w.Code)
		}
	})
}
