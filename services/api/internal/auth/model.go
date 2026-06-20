package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===== MongoDB documents =====

// Shop — ร้านค้า (1 user = 1 shop ใน Phase 1)
type Shop struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `bson:"owner_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
}

// User — เจ้าของร้าน เก็บเฉพาะ password_hash (ไม่เก็บ plain text)
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	ShopID       primitive.ObjectID `bson:"shop_id"`
	CreatedAt    time.Time          `bson:"created_at"`
}

// ===== Request / Response DTO (snake_case ตาม API Contract) =====

type RegisterRequest struct {
	ShopName string `json:"shop_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	Shop  ShopInfo `json:"shop"`
}

type ShopInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
