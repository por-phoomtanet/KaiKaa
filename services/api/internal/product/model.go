package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product — สินค้าใน 1 ร้าน (ผูกกับ shop_id เสมอ)
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ShopID    primitive.ObjectID `bson:"shop_id"`
	Name      string             `bson:"name"`
	Price     int                `bson:"price"` // บาท (integer)
	Emoji     string             `bson:"emoji"`
	IsActive  bool               `bson:"is_active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ===== DTO (snake_case ตาม API Contract) =====

// ProductRequest ใช้ทั้ง create และ update
type ProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"gte=0"`
	Emoji string `json:"emoji"`
}

// ProductResponse — รูปแบบที่ส่งกลับ (ตัด shop_id/timestamp ออกตาม Contract)
type ProductResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Emoji string `json:"emoji"`
}

// DeleteResponse ตาม Contract: { "id": "...", "deleted": true }
type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func toResponse(p Product) ProductResponse {
	return ProductResponse{
		ID:    p.ID.Hex(),
		Name:  p.Name,
		Price: p.Price,
		Emoji: p.Emoji,
	}
}
