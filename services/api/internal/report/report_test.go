package report

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kaikaa/api/internal/sale"
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

func insertSale(t *testing.T, db *mongo.Database, shop, product primitive.ObjectID, name, emoji string, total int, method string, soldAt time.Time) {
	t.Helper()
	_, err := db.Collection("sales").InsertOne(ctx, sale.Sale{
		ID: primitive.NewObjectID(), ShopID: shop, ProductID: product,
		Name: name, Emoji: emoji, Total: total, Method: method, SoldAt: soldAt,
	})
	if err != nil {
		t.Fatalf("insert sale: %v", err)
	}
}

// noonICT คืนเวลา 12:00 ICT ของวันที่ระบุ (กันคร่อมเส้นเที่ยงคืน)
func noonICT(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 12, 0, 0, 0, ictLoc)
}

func TestDailyAggregation(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	pa := primitive.NewObjectID()
	pb := primitive.NewObjectID()

	now := time.Now().In(ictLoc)
	at := noonICT(now.Year(), now.Month(), now.Day())
	insertSale(t, db, shop, pa, "A", "☕", 50, "cash", at)
	insertSale(t, db, shop, pa, "A", "☕", 50, "cash", at)
	insertSale(t, db, shop, pb, "B", "🧋", 100, "transfer", at)

	rep, err := s.Daily(ctx, shop, at.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("daily: %v", err)
	}
	if rep.Total != 200 || rep.Count != 3 {
		t.Errorf("total=%d count=%d, want 200/3", rep.Total, rep.Count)
	}
	if rep.Cash != 100 || rep.Transfer != 100 {
		t.Errorf("cash=%d transfer=%d, want 100/100", rep.Cash, rep.Transfer)
	}
	if rep.CashPct+rep.TransferPct != 100 {
		t.Errorf("cash_pct + transfer_pct = %d, want 100", rep.CashPct+rep.TransferPct)
	}
	if len(rep.WeekBars) != 7 {
		t.Errorf("week_bars len = %d, want 7", len(rep.WeekBars))
	}
	if last := rep.WeekBars[6]; last.Amount != 200 {
		t.Errorf("last week bar amount = %d, want 200", last.Amount)
	}
	if len(rep.BestSellers) != 2 || rep.BestSellers[0].Name != "A" || rep.BestSellers[0].Qty != 2 {
		t.Errorf("best sellers wrong: %+v", rep.BestSellers)
	}
	if rep.BestSellers[0].PctWidth != 100 || rep.BestSellers[0].Rank != 1 {
		t.Errorf("top seller pct/rank wrong: %+v", rep.BestSellers[0])
	}
	if len(rep.Recent) != 3 {
		t.Errorf("recent len = %d, want 3", len(rep.Recent))
	}
}

func TestDailyEmptyReturnsZeros(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()

	rep, err := s.Daily(ctx, shop, "2020-01-01")
	if err != nil {
		t.Fatalf("daily empty should not error: %v", err)
	}
	if rep.Total != 0 || rep.Count != 0 || rep.CashPct != 0 || rep.TransferPct != 0 {
		t.Errorf("empty day should be zeros: %+v", rep)
	}
	if len(rep.WeekBars) != 7 {
		t.Errorf("week_bars len = %d, want 7 even when empty", len(rep.WeekBars))
	}
	if len(rep.BestSellers) != 0 || len(rep.Recent) != 0 {
		t.Errorf("best_sellers/recent should be empty slices")
	}
}

func TestDailyRecentCappedAt6(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	now := time.Now().In(ictLoc)
	base := noonICT(now.Year(), now.Month(), now.Day())
	for i := 0; i < 8; i++ {
		insertSale(t, db, shop, primitive.NewObjectID(), "x", "x", 10, "cash", base.Add(time.Duration(i)*time.Minute))
	}
	rep, _ := s.Daily(ctx, shop, base.Format("2006-01-02"))
	if len(rep.Recent) != 6 {
		t.Errorf("recent len = %d, want capped at 6", len(rep.Recent))
	}
}

func TestMonthlyBucketsByWeek(t *testing.T) {
	db := newTestDB(t)
	s := NewService(db)
	shop := primitive.NewObjectID()
	p := primitive.NewObjectID()

	insertSale(t, db, shop, p, "x", "x", 100, "cash", noonICT(2026, time.March, 3))  // week 1
	insertSale(t, db, shop, p, "x", "x", 200, "cash", noonICT(2026, time.March, 10)) // week 2
	insertSale(t, db, shop, p, "x", "x", 300, "cash", noonICT(2026, time.March, 20)) // week 3
	insertSale(t, db, shop, p, "x", "x", 400, "cash", noonICT(2026, time.March, 25)) // week 4

	rep, err := s.Monthly(ctx, shop, "2026-03")
	if err != nil {
		t.Fatalf("monthly: %v", err)
	}
	if rep.MonthTotal != 1000 || rep.BillCount != 4 {
		t.Errorf("total=%d count=%d, want 1000/4", rep.MonthTotal, rep.BillCount)
	}
	if len(rep.WeekBars) != 4 {
		t.Fatalf("week_bars len = %d, want 4", len(rep.WeekBars))
	}
	want := []int{100, 200, 300, 400}
	for i, w := range want {
		if rep.WeekBars[i].Amount != w {
			t.Errorf("week %d amount = %d, want %d", i+1, rep.WeekBars[i].Amount, w)
		}
	}
	if rep.WeekBars[0].Label != "สัปดาห์ 1" {
		t.Errorf("label = %q, want 'สัปดาห์ 1'", rep.WeekBars[0].Label)
	}
	if rep.AvgPerBill != 250 {
		t.Errorf("avg_per_bill = %d, want 250", rep.AvgPerBill)
	}
	if rep.AvgPerDay != 32 { // round(1000/31)
		t.Errorf("avg_per_day = %d, want 32", rep.AvgPerDay)
	}
}

func TestMonthlyEmpty(t *testing.T) {
	s := NewService(newTestDB(t))
	rep, err := s.Monthly(ctx, primitive.NewObjectID(), "2020-01")
	if err != nil {
		t.Fatalf("monthly empty should not error: %v", err)
	}
	if rep.MonthTotal != 0 || rep.BillCount != 0 || rep.AvgPerBill != 0 {
		t.Errorf("empty month should be zeros: %+v", rep)
	}
	if len(rep.WeekBars) != 4 {
		t.Errorf("week_bars len = %d, want 4", len(rep.WeekBars))
	}
}

func TestInvalidDateMonth(t *testing.T) {
	s := NewService(newTestDB(t))
	shop := primitive.NewObjectID()
	if _, err := s.Daily(ctx, shop, "2026/06/21"); err != ErrInvalidDate {
		t.Errorf("daily err = %v, want ErrInvalidDate", err)
	}
	if _, err := s.Monthly(ctx, shop, "June"); err != ErrInvalidMonth {
		t.Errorf("monthly err = %v, want ErrInvalidMonth", err)
	}
}
