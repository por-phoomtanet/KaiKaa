package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// อายุ token — 7 วัน (Phase 1 เก็บใน AsyncStorage บนมือถือ)
const ttl = 7 * 24 * time.Hour

// Claims ที่ฝังใน JWT — shop_id ใช้ filter ข้อมูลแบบ multi-tenant
type Claims struct {
	UserID string `json:"user_id"`
	ShopID string `json:"shop_id"`
	jwt.RegisteredClaims
}

// Generate สร้าง signed JWT (HS256)
func Generate(secret, userID, shopID string) (string, error) {
	claims := Claims{
		UserID: userID,
		ShopID: shopID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// Parse ตรวจ signature + อายุ แล้วคืน claims
func Parse(secret, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
