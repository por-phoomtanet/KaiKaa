package product

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

var ctx = context.Background()

func TestCreateAndList(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()

	created, err := s.Create(ctx, shop, ProductRequest{Name: "ลาเต้", Price: 55, Emoji: "☕"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" || created.Name != "ลาเต้" || created.Price != 55 {
		t.Errorf("unexpected created: %+v", created)
	}

	list, err := s.List(ctx, shop)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Errorf("list = %+v, want 1 item matching created", list)
	}
}

func TestListIsolatedPerShop(t *testing.T) {
	s := NewService(newTestDB(t))
	shopA := primitive.NewObjectID()
	shopB := primitive.NewObjectID()

	_, _ = s.Create(ctx, shopA, ProductRequest{Name: "A1", Price: 10})
	_, _ = s.Create(ctx, shopA, ProductRequest{Name: "A2", Price: 20})
	_, _ = s.Create(ctx, shopB, ProductRequest{Name: "B1", Price: 30})

	listA, _ := s.List(ctx, shopA)
	if len(listA) != 2 {
		t.Errorf("shopA list len = %d, want 2", len(listA))
	}
	listB, _ := s.List(ctx, shopB)
	if len(listB) != 1 || listB[0].Name != "B1" {
		t.Errorf("shopB list = %+v, want only B1", listB)
	}
}

func TestUpdate(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()
	created, _ := s.Create(ctx, shop, ProductRequest{Name: "old", Price: 10, Emoji: "☕"})

	updated, err := s.Update(ctx, shop, created.ID, ProductRequest{Name: "new", Price: 99, Emoji: "🧋"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "new" || updated.Price != 99 || updated.Emoji != "🧋" {
		t.Errorf("unexpected updated: %+v", updated)
	}
	if updated.ID != created.ID {
		t.Errorf("id changed: %s -> %s", created.ID, updated.ID)
	}
}

func TestUpdateOtherShopReturnsNotFound(t *testing.T) {
	s := NewService(newTestDB(t))
	owner := primitive.NewObjectID()
	attacker := primitive.NewObjectID()
	created, _ := s.Create(ctx, owner, ProductRequest{Name: "x", Price: 10})

	if _, err := s.Update(ctx, attacker, created.ID, ProductRequest{Name: "hacked", Price: 0}); err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()
	created, _ := s.Create(ctx, shop, ProductRequest{Name: "x", Price: 10})

	if err := s.Delete(ctx, shop, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	list, _ := s.List(ctx, shop)
	if len(list) != 0 {
		t.Errorf("after delete list len = %d, want 0", len(list))
	}
}

func TestDeleteOtherShopReturnsNotFound(t *testing.T) {
	s := NewService(newTestDB(t))
	owner := primitive.NewObjectID()
	attacker := primitive.NewObjectID()
	created, _ := s.Create(ctx, owner, ProductRequest{Name: "x", Price: 10})

	if err := s.Delete(ctx, attacker, created.ID); err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
	// ของเจ้าของต้องยังอยู่
	list, _ := s.List(ctx, owner)
	if len(list) != 1 {
		t.Errorf("owner product was deleted by attacker; list len = %d", len(list))
	}
}

func TestUpdateDeleteInvalidID(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()

	if _, err := s.Update(ctx, shop, "not-a-valid-id", ProductRequest{Name: "x"}); err != ErrNotFound {
		t.Errorf("update invalid id err = %v, want ErrNotFound", err)
	}
	if err := s.Delete(ctx, shop, "not-a-valid-id"); err != ErrNotFound {
		t.Errorf("delete invalid id err = %v, want ErrNotFound", err)
	}
}
