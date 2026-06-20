package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/internal/response"
	"github.com/kaikaa/api/internal/token"
)

// keys ที่เก็บใน gin.Context ให้ handler ปลายทางดึงไปใช้
const (
	CtxUserID = "user_id"
	CtxShopID = "shop_id"
)

// JWT ตรวจ Authorization: Bearer <token>
// ผ่านแล้ว set user_id / shop_id ลง context สำหรับ filter แบบ multi-tenant
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := token.Parse(secret, tokenStr)
		if err != nil {
			response.Unauthorized(c)
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxShopID, claims.ShopID)
		c.Next()
	}
}
