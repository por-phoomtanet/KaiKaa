package report

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/kaikaa/api/internal/sale"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidDate  = errors.New("invalid date format, want YYYY-MM-DD")
	ErrInvalidMonth = errors.New("invalid month format, want YYYY-MM")
)

type Service struct {
	sales *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{sales: db.Collection("sales")}
}

// fetchRange ดึง sales ของ shop ในช่วงเวลา [start, end)
func (s *Service) fetchRange(ctx context.Context, shopID primitive.ObjectID, start, end time.Time) ([]sale.Sale, error) {
	cur, err := s.sales.Find(ctx, bson.M{
		"shop_id": shopID,
		"sold_at": bson.M{"$gte": start, "$lt": end},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []sale.Sale
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func roundPct(part, whole int) int {
	if whole == 0 {
		return 0
	}
	return int(math.Round(float64(part) / float64(whole) * 100))
}

// Daily สร้างรายงานรายวันของวันที่ระบุ (ICT) — date ว่าง = วันนี้
func (s *Service) Daily(ctx context.Context, shopID primitive.ObjectID, dateStr string) (DailyReport, error) {
	if dateStr == "" {
		dateStr = time.Now().In(ictLoc).Format("2006-01-02")
	}
	day, err := time.ParseInLocation("2006-01-02", dateStr, ictLoc)
	if err != nil {
		return DailyReport{}, ErrInvalidDate
	}

	// ดึง 7 วันย้อนหลังถึงวันที่เลือก (ใช้ทั้ง week_bars และคำนวณ change vs ค่าเฉลี่ย)
	rangeStart := day.AddDate(0, 0, -6).UTC()
	rangeEnd := day.AddDate(0, 0, 1).UTC()
	sales, err := s.fetchRange(ctx, shopID, rangeStart, rangeEnd)
	if err != nil {
		return DailyReport{}, err
	}

	targetKey := day.Format("2006-01-02")
	dayTotals := map[string]int{}
	var today []sale.Sale
	for _, sl := range sales {
		key := sl.SoldAt.In(ictLoc).Format("2006-01-02")
		dayTotals[key] += sl.Total
		if key == targetKey {
			today = append(today, sl)
		}
	}

	rep := DailyReport{
		BestSellers: []BestSeller{},
		WeekBars:    make([]DayBar, 0, 7),
		Recent:      []RecentSale{},
	}

	// week_bars: เก่า→ใหม่ (ตัวสุดท้าย = วันที่เลือก) เติม 0 วันที่ไม่มีขาย
	prevSum, prevDays := 0, 0
	for i := 6; i >= 0; i-- {
		d := day.AddDate(0, 0, -i)
		amt := dayTotals[d.Format("2006-01-02")]
		rep.WeekBars = append(rep.WeekBars, DayBar{Day: thaiWeekday(d.Weekday()), Amount: amt})
		if i != 0 { // วันก่อนหน้า (ไม่รวมวันนี้) ใช้หาค่าเฉลี่ย
			prevSum += amt
			prevDays++
		}
	}

	// total / cash / transfer / count ของวันที่เลือก
	for _, sl := range today {
		rep.Total += sl.Total
		rep.Count++
		if sl.Method == "cash" {
			rep.Cash += sl.Total
		} else {
			rep.Transfer += sl.Total
		}
	}
	rep.CashPct = roundPct(rep.Cash, rep.Total)
	if rep.Total > 0 {
		rep.TransferPct = 100 - rep.CashPct // การันตี cash_pct + transfer_pct = 100
	}

	// change_pct: เทียบยอดวันนี้กับค่าเฉลี่ย 6 วันก่อนหน้า
	if prevDays > 0 {
		avg := float64(prevSum) / float64(prevDays)
		if avg > 0 {
			rep.ChangePct = int(math.Round((float64(rep.Total) - avg) / avg * 100))
		}
	}

	rep.BestSellers = bestSellers(today)
	rep.Recent = recentSales(today)
	return rep, nil
}

// bestSellers จัดอันดับสินค้าขายดีของวัน (Top 4) ตามจำนวนแก้ว
func bestSellers(sales []sale.Sale) []BestSeller {
	type agg struct {
		name, emoji string
		qty, rev    int
	}
	m := map[primitive.ObjectID]*agg{}
	order := []primitive.ObjectID{}
	for _, sl := range sales {
		a, ok := m[sl.ProductID]
		if !ok {
			a = &agg{name: sl.Name, emoji: sl.Emoji}
			m[sl.ProductID] = a
			order = append(order, sl.ProductID)
		}
		a.qty++
		a.rev += sl.Total
	}

	list := make([]*agg, 0, len(order))
	for _, id := range order {
		list = append(list, m[id])
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].qty > list[j].qty })

	out := []BestSeller{}
	maxQty := 0
	if len(list) > 0 {
		maxQty = list[0].qty
	}
	for i, a := range list {
		if i >= 4 {
			break
		}
		out = append(out, BestSeller{
			Rank: i + 1, Name: a.name, Emoji: a.emoji,
			Qty: a.qty, Revenue: a.rev, PctWidth: roundPct(a.qty, maxQty),
		})
	}
	return out
}

// recentSales คืนรายการล่าสุด 6 รายการ (ใหม่ก่อน)
func recentSales(sales []sale.Sale) []RecentSale {
	sorted := make([]sale.Sale, len(sales))
	copy(sorted, sales)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].SoldAt.After(sorted[j].SoldAt) })

	out := []RecentSale{}
	for i, sl := range sorted {
		if i >= 6 {
			break
		}
		out = append(out, RecentSale{
			Name: sl.Name, Emoji: sl.Emoji,
			Time:   sl.SoldAt.In(ictLoc).Format("15:04"),
			Total:  sl.Total,
			Method: sl.Method,
		})
	}
	return out
}

// Monthly สร้างรายงานรายเดือน — month ว่าง = เดือนนี้
func (s *Service) Monthly(ctx context.Context, shopID primitive.ObjectID, monthStr string) (MonthlyReport, error) {
	if monthStr == "" {
		monthStr = time.Now().In(ictLoc).Format("2006-01")
	}
	first, err := time.ParseInLocation("2006-01", monthStr, ictLoc)
	if err != nil {
		return MonthlyReport{}, ErrInvalidMonth
	}
	next := first.AddDate(0, 1, 0)

	sales, err := s.fetchRange(ctx, shopID, first.UTC(), next.UTC())
	if err != nil {
		return MonthlyReport{}, err
	}

	rep := MonthlyReport{WeekBars: make([]WeekBar, 4)}
	weeks := [4]int{}
	for _, sl := range sales {
		rep.MonthTotal += sl.Total
		rep.BillCount++
		// แบ่งสัปดาห์ตามวันที่ในเดือน: 1-7, 8-14, 15-21, 22-สิ้นเดือน
		idx := (sl.SoldAt.In(ictLoc).Day() - 1) / 7
		if idx > 3 {
			idx = 3
		}
		weeks[idx] += sl.Total
	}
	for i := 0; i < 4; i++ {
		rep.WeekBars[i] = WeekBar{Label: "สัปดาห์ " + string(rune('1'+i)), Amount: weeks[i]}
	}

	daysInMonth := int(next.Sub(first).Hours() / 24)
	if daysInMonth > 0 {
		rep.AvgPerDay = int(math.Round(float64(rep.MonthTotal) / float64(daysInMonth)))
	}
	if rep.BillCount > 0 {
		rep.AvgPerBill = int(math.Round(float64(rep.MonthTotal) / float64(rep.BillCount)))
	}
	return rep, nil
}
