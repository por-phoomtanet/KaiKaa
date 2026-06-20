package sale

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ICT — เวลาไทย (UTC+7, ไม่มี DST) ใช้ FixedZone จึงไม่ต้องพึ่ง tzdata ใน container
var ictLoc = time.FixedZone("ICT", 7*3600)

// Sale — 1 รายการขาย เก็บ snapshot ชื่อ/ราคา/emoji ของสินค้า ณ เวลาที่ขาย
// (ไม่อ้างอิงสด เพราะสินค้าอาจถูกแก้ราคา/ลบภายหลัง แต่บิลเก่าต้องคงเดิม)
type Sale struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ShopID    primitive.ObjectID `bson:"shop_id"`
	ProductID primitive.ObjectID `bson:"product_id"`
	Name      string             `bson:"product_name"`
	Emoji     string             `bson:"product_emoji"`
	Total     int                `bson:"total"`
	Method    string             `bson:"method"` // cash | transfer
	SoldAt    time.Time          `bson:"sold_at"`
}

// ===== DTO (snake_case ตาม API Contract) =====

type CreateSaleRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Method    string `json:"method" binding:"required,oneof=cash transfer"`
}

type SaleResponse struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji"`
	Total     int    `json:"total"`
	Method    string `json:"method"`
	Time      string `json:"time"`    // HH:mm (เวลาไทย)
	SoldAt    string `json:"sold_at"` // RFC3339 (UTC)
}

func toResponse(s Sale) SaleResponse {
	return SaleResponse{
		ID:        s.ID.Hex(),
		ProductID: s.ProductID.Hex(),
		Name:      s.Name,
		Emoji:     s.Emoji,
		Total:     s.Total,
		Method:    s.Method,
		Time:      s.SoldAt.In(ictLoc).Format("15:04"),
		SoldAt:    s.SoldAt.UTC().Format(time.RFC3339),
	}
}
