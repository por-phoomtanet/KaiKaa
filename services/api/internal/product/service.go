package product

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrNotFound — สินค้าไม่มี หรือไม่ใช่ของ shop นี้ (handler map เป็น 404)
var ErrNotFound = errors.New("product not found")

type Service struct {
	products *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{products: db.Collection("products")}
}

// List คืนสินค้าทั้งหมดของ shop (เรียงตามเวลาสร้าง)
func (s *Service) List(ctx context.Context, shopID primitive.ObjectID) ([]ProductResponse, error) {
	cur, err := s.products.Find(ctx,
		bson.M{"shop_id": shopID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []Product
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}

	out := make([]ProductResponse, 0, len(docs))
	for _, p := range docs {
		out = append(out, toResponse(p))
	}
	return out, nil
}

// Create เพิ่มสินค้าใหม่ผูกกับ shop
func (s *Service) Create(ctx context.Context, shopID primitive.ObjectID, req ProductRequest) (ProductResponse, error) {
	now := time.Now()
	p := Product{
		ID:        primitive.NewObjectID(),
		ShopID:    shopID,
		Name:      req.Name,
		Price:     req.Price,
		Emoji:     req.Emoji,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := s.products.InsertOne(ctx, p); err != nil {
		return ProductResponse{}, err
	}
	return toResponse(p), nil
}

// Update แก้ไขสินค้า — filter ด้วย shop_id ด้วย กันแก้ของร้านอื่น
func (s *Service) Update(ctx context.Context, shopID primitive.ObjectID, idHex string, req ProductRequest) (ProductResponse, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ProductResponse{}, ErrNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":       req.Name,
		"price":      req.Price,
		"emoji":      req.Emoji,
		"updated_at": time.Now(),
	}}

	var updated Product
	err = s.products.FindOneAndUpdate(ctx,
		bson.M{"_id": id, "shop_id": shopID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ProductResponse{}, ErrNotFound
		}
		return ProductResponse{}, err
	}
	return toResponse(updated), nil
}

// Delete ลบสินค้า — filter ด้วย shop_id ด้วย กันลบของร้านอื่น
func (s *Service) Delete(ctx context.Context, shopID primitive.ObjectID, idHex string) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}

	res, err := s.products.DeleteOne(ctx, bson.M{"_id": id, "shop_id": shopID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
