package sale

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// ErrProductNotFound — product_id ไม่มี หรือไม่ใช่ของ shop นี้
	ErrProductNotFound = errors.New("product not found")
	// ErrInvalidDate — รูปแบบวันที่ผิด ต้องเป็น YYYY-MM-DD
	ErrInvalidDate = errors.New("invalid date format, want YYYY-MM-DD")
)

type Service struct {
	sales    *mongo.Collection
	products *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		sales:    db.Collection("sales"),
		products: db.Collection("products"),
	}
}

// Create บันทึกการขาย — snapshot ชื่อ/ราคา/emoji จากสินค้าจริงของ shop
func (s *Service) Create(ctx context.Context, shopID primitive.ObjectID, req CreateSaleRequest) (SaleResponse, error) {
	pid, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return SaleResponse{}, ErrProductNotFound
	}

	// ดึงสินค้าที่เป็นของ shop นี้เท่านั้น (กันขายสินค้าร้านอื่น)
	var prod struct {
		Name  string `bson:"name"`
		Price int    `bson:"price"`
		Emoji string `bson:"emoji"`
	}
	err = s.products.FindOne(ctx, bson.M{"_id": pid, "shop_id": shopID}).Decode(&prod)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return SaleResponse{}, ErrProductNotFound
		}
		return SaleResponse{}, err
	}

	sale := Sale{
		ID:        primitive.NewObjectID(),
		ShopID:    shopID,
		ProductID: pid,
		Name:      prod.Name,
		Emoji:     prod.Emoji,
		Total:     prod.Price, // ราคา server เป็นคนกำหนด ไม่รับจาก client
		Method:    req.Method,
		SoldAt:    time.Now(),
	}
	if _, err := s.sales.InsertOne(ctx, sale); err != nil {
		return SaleResponse{}, err
	}
	return toResponse(sale), nil
}

// ListByDate คืนรายการขายของวันที่ระบุ (ตามเวลาไทย) เรียงใหม่ก่อน
// date ว่าง = วันนี้
func (s *Service) ListByDate(ctx context.Context, shopID primitive.ObjectID, dateStr string) ([]SaleResponse, error) {
	if dateStr == "" {
		dateStr = time.Now().In(ictLoc).Format("2006-01-02")
	}
	day, err := time.ParseInLocation("2006-01-02", dateStr, ictLoc)
	if err != nil {
		return nil, ErrInvalidDate
	}
	start := day.UTC()
	end := day.Add(24 * time.Hour).UTC()

	cur, err := s.sales.Find(ctx,
		bson.M{
			"shop_id": shopID,
			"sold_at": bson.M{"$gte": start, "$lt": end},
		},
		options.Find().SetSort(bson.D{{Key: "sold_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []Sale
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	out := make([]SaleResponse, 0, len(docs))
	for _, d := range docs {
		out = append(out, toResponse(d))
	}
	return out, nil
}
