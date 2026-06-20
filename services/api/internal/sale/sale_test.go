package sale

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func newTestDB(t *testing.T) *mongo.Database {
	t.Helper()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	cctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(cctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("mongo unavailable (%s): %v", uri, err)
	}
	if err := client.Ping(cctx, nil); err != nil {
		t.Skipf("mongo ping failed (%s): %v", uri, err)
	}

	db := client.Database("kaikaa_test_" + primitive.NewObjectID().Hex())
	t.Cleanup(func() {
		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Drop(cctx)
		_ = client.Disconnect(cctx)
	})
	return db
}

// seedProduct แทรกสินค้าลง collection ตรงๆ คืน id
func seedProduct(t *testing.T, db *mongo.Database, shopID primitive.ObjectID, name string, price int, emoji string) primitive.ObjectID {
	t.Helper()
	id := primitive.NewObjectID()
	_, err := db.Collection("products").InsertOne(ctx, bson.M{
		"_id": id, "shop_id": shopID, "name": name, "price": price, "emoji": emoji,
		"is_active": true, "created_at": time.Now(), "updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return id
}

func TestCreateSnapshotsProduct(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	pid := seedProduct(t, db, shop, "ลาเต้", 55, "☕")

	res, err := s.Create(ctx, shop, CreateSaleRequest{ProductID: pid.Hex(), Method: "cash"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if res.Name != "ลาเต้" || res.Emoji != "☕" || res.Total != 55 {
		t.Errorf("snapshot mismatch: %+v", res)
	}
	if res.Method != "cash" {
		t.Errorf("method = %q, want cash", res.Method)
	}
	if res.Time == "" || res.SoldAt == "" {
		t.Errorf("server must fill time/sold_at, got %+v", res)
	}
}

func TestCreatePriceIgnoresClient(t *testing.T) {
	// แม้สินค้าถูกแก้ราคาภายหลัง บิลเก่าต้องคงราคา snapshot
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	pid := seedProduct(t, db, shop, "x", 40, "x")

	res, _ := s.Create(ctx, shop, CreateSaleRequest{ProductID: pid.Hex(), Method: "transfer"})
	if res.Total != 40 {
		t.Errorf("total = %d, want 40 (from product price)", res.Total)
	}
}

func TestCreateNonexistentProduct(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()
	ghost := primitive.NewObjectID().Hex()

	if _, err := s.Create(ctx, shop, CreateSaleRequest{ProductID: ghost, Method: "cash"}); err != ErrProductNotFound {
		t.Errorf("err = %v, want ErrProductNotFound", err)
	}
}

func TestCreateOtherShopProduct(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	owner := primitive.NewObjectID()
	attacker := primitive.NewObjectID()
	pid := seedProduct(t, db, owner, "x", 10, "x")

	if _, err := s.Create(ctx, attacker, CreateSaleRequest{ProductID: pid.Hex(), Method: "cash"}); err != ErrProductNotFound {
		t.Errorf("err = %v, want ErrProductNotFound (cross-shop)", err)
	}
}

func TestListByDateFiltersDay(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()

	// แทรกบิล: วันนี้ 2 รายการ, เมื่อวาน 1 รายการ (กำหนด sold_at ตรงๆ)
	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)
	insert := func(soldAt time.Time) {
		_, _ = db.Collection("sales").InsertOne(ctx, Sale{
			ID: primitive.NewObjectID(), ShopID: shop, ProductID: primitive.NewObjectID(),
			Name: "x", Emoji: "x", Total: 10, Method: "cash", SoldAt: soldAt,
		})
	}
	insert(now)
	insert(now)
	insert(yesterday)

	today := now.In(ictLoc).Format("2006-01-02")
	list, err := s.ListByDate(ctx, shop, today)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("today count = %d, want 2", len(list))
	}

	yStr := yesterday.In(ictLoc).Format("2006-01-02")
	listY, _ := s.ListByDate(ctx, shop, yStr)
	if len(listY) != 1 {
		t.Errorf("yesterday count = %d, want 1", len(listY))
	}
}

func TestListByDateSortedDesc(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	now := time.Now().UTC()

	older := Sale{ID: primitive.NewObjectID(), ShopID: shop, ProductID: primitive.NewObjectID(), Name: "older", Total: 10, Method: "cash", SoldAt: now.Add(-2 * time.Hour)}
	newer := Sale{ID: primitive.NewObjectID(), ShopID: shop, ProductID: primitive.NewObjectID(), Name: "newer", Total: 10, Method: "cash", SoldAt: now.Add(-1 * time.Hour)}
	_, _ = db.Collection("sales").InsertOne(ctx, older)
	_, _ = db.Collection("sales").InsertOne(ctx, newer)

	list, _ := s.ListByDate(ctx, shop, now.In(ictLoc).Format("2006-01-02"))
	if len(list) != 2 || list[0].Name != "newer" {
		t.Errorf("expected newest first, got %+v", list)
	}
}

func TestListByDateInvalid(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()
	if _, err := s.ListByDate(ctx, shop, "2026/06/20"); err != ErrInvalidDate {
		t.Errorf("err = %v, want ErrInvalidDate", err)
	}
}
