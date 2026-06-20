package ai

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/kaikaa/api/internal/report"
	"github.com/kaikaa/api/internal/sale"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// mockLLM — Completer ปลอม ใช้ test โดยไม่ต้องมี OPENROUTER_KEY/network
type mockLLM struct {
	reply string
	err   error
}

func (m mockLLM) Complete(_ context.Context, _ string) (string, error) {
	return m.reply, m.err
}

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
		t.Skipf("mongo unavailable: %v", err)
	}
	if err := client.Ping(cctx, nil); err != nil {
		t.Skipf("mongo ping failed: %v", err)
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

func seedToday(t *testing.T, db *mongo.Database, shop primitive.ObjectID) {
	t.Helper()
	loc := time.FixedZone("ICT", 7*3600)
	td := time.Now().In(loc)
	noon := time.Date(td.Year(), td.Month(), td.Day(), 12, 0, 0, 0, loc)
	pid := primitive.NewObjectID()
	for i := 0; i < 2; i++ {
		_, err := db.Collection("sales").InsertOne(ctx, sale.Sale{
			ID: primitive.NewObjectID(), ShopID: shop, ProductID: pid,
			Name: "ลาเต้", Emoji: "☕", Total: 55, Method: "cash", SoldAt: noon,
		})
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
}

// ===== pure unit: parseAIJSON =====

func TestParseAIPlainJSON(t *testing.T) {
	p, err := parseAIJSON(`{"summary":"ดีมาก","recommendations":["a","b","c"]}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if p.Summary != "ดีมาก" || len(p.Recommendations) != 3 {
		t.Errorf("unexpected: %+v", p)
	}
}

func TestParseAIWithMarkdownFence(t *testing.T) {
	raw := "นี่คือผลลัพธ์:\n```json\n{\"summary\":\"x\",\"recommendations\":[\"a\"]}\n```\nจบ"
	p, err := parseAIJSON(raw)
	if err != nil {
		t.Fatalf("parse fenced: %v", err)
	}
	if p.Summary != "x" || len(p.Recommendations) != 1 {
		t.Errorf("unexpected: %+v", p)
	}
}

func TestParseAINoJSON(t *testing.T) {
	if _, err := parseAIJSON("ขอโทษ ฉันไม่เข้าใจ"); err == nil {
		t.Error("expected error when no JSON present")
	}
}

// ===== service with mock LLM (needs mongo for report data) =====

func TestSummaryHappyPath(t *testing.T) {
	db := newTestDB(t)
	shop := primitive.NewObjectID()
	seedToday(t, db, shop)

	llm := mockLLM{reply: `{"summary":"วันนี้ขายดี","recommendations":["เพิ่มสต็อกลาเต้","ออกโปรช่วงบ่าย","กระตุ้นการโอน"]}`}
	s := NewService(report.NewService(db), llm)

	res, err := s.Summary(ctx, shop, "")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if res.Summary != "วันนี้ขายดี" {
		t.Errorf("summary = %q", res.Summary)
	}
	if len(res.Recommendations) != 3 {
		t.Errorf("recommendations len = %d, want 3", len(res.Recommendations))
	}
	// best_sellers ต้องมาจากข้อมูลจริง (ลาเต้ 2 แก้ว) ไม่ใช่จาก AI
	if len(res.BestSellers) != 1 || res.BestSellers[0].Name != "ลาเต้" || res.BestSellers[0].Qty != 2 {
		t.Errorf("best_sellers from data wrong: %+v", res.BestSellers)
	}
}

func TestSummaryLLMErrorPropagates(t *testing.T) {
	db := newTestDB(t)
	shop := primitive.NewObjectID()
	seedToday(t, db, shop)

	s := NewService(report.NewService(db), mockLLM{err: errors.New("timeout")})
	if _, err := s.Summary(ctx, shop, ""); err == nil {
		t.Error("expected error when LLM fails, got nil")
	}
}

func TestSummaryInvalidDate(t *testing.T) {
	db := newTestDB(t)
	s := NewService(report.NewService(db), mockLLM{reply: `{"summary":"x","recommendations":[]}`})
	_, err := s.Summary(ctx, primitive.NewObjectID(), "bad-date")
	if !errors.Is(err, report.ErrInvalidDate) {
		t.Errorf("err = %v, want ErrInvalidDate", err)
	}
}
