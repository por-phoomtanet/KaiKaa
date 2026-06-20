package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/internal/token"
)

const testSecret = "test-secret"

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(JWT(testSecret))
	r.GET("/p", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.GetString(CtxUserID),
			"shop_id": c.GetString(CtxShopID),
		})
	})
	return r
}

func do(r *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestNoTokenReturns401(t *testing.T) {
	w := do(testRouter(), "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code = %d, want 401", w.Code)
	}
}

func TestMalformedHeaderReturns401(t *testing.T) {
	w := do(testRouter(), "token-without-bearer-prefix")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code = %d, want 401", w.Code)
	}
}

func TestInvalidTokenReturns401(t *testing.T) {
	w := do(testRouter(), "Bearer not-a-valid-jwt")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("code = %d, want 401", w.Code)
	}
}

func TestValidTokenReturns200AndSetsContext(t *testing.T) {
	tok, _ := token.Generate(testSecret, "user1", "shop1")
	w := do(testRouter(), "Bearer "+tok)
	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", w.Code)
	}
	body := w.Body.String()
	if !contains(body, "shop1") || !contains(body, "user1") {
		t.Errorf("body %q missing user_id/shop_id from context", body)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
