package report

import "time"

// ICT — เวลาไทย (UTC+7) ใช้ตัดวัน/เดือน ให้ตรงกับ sale package
var ictLoc = time.FixedZone("ICT", 7*3600)

// ===== Daily report (ตาม API Contract) =====

type DailyReport struct {
	Total       int          `json:"total"`
	Count       int          `json:"count"`
	Cash        int          `json:"cash"`
	Transfer    int          `json:"transfer"`
	CashPct     int          `json:"cash_pct"`
	TransferPct int          `json:"transfer_pct"`
	ChangePct   int          `json:"change_pct"`
	BestSellers []BestSeller `json:"best_sellers"`
	WeekBars    []DayBar     `json:"week_bars"`
	Recent      []RecentSale `json:"recent"`
}

type BestSeller struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	Emoji    string `json:"emoji"`
	Qty      int    `json:"qty"`
	Revenue  int    `json:"revenue"`
	PctWidth int    `json:"pct_width"`
}

type DayBar struct {
	Day    string `json:"day"`
	Amount int    `json:"amount"`
}

type RecentSale struct {
	Name   string `json:"name"`
	Emoji  string `json:"emoji"`
	Time   string `json:"time"`
	Total  int    `json:"total"`
	Method string `json:"method"`
}

// ===== Monthly report (ตาม API Contract) =====

type MonthlyReport struct {
	MonthTotal int       `json:"month_total"`
	WeekBars   []WeekBar `json:"week_bars"`
	AvgPerDay  int       `json:"avg_per_day"`
	BillCount  int       `json:"bill_count"`
	AvgPerBill int       `json:"avg_per_bill"`
}

type WeekBar struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// thaiWeekday ย่อชื่อวันภาษาไทยสำหรับ bar chart รายวัน
func thaiWeekday(w time.Weekday) string {
	switch w {
	case time.Monday:
		return "จ"
	case time.Tuesday:
		return "อ"
	case time.Wednesday:
		return "พ"
	case time.Thursday:
		return "พฤ"
	case time.Friday:
		return "ศ"
	case time.Saturday:
		return "ส"
	default:
		return "อา"
	}
}
