package ai

// SummaryRequest — body ว่าง {} = วันนี้ หรือระบุ date
type SummaryRequest struct {
	Date string `json:"date"`
}

// BestSeller — ข้อมูลจริงจากยอดขาย (ไม่ได้มาจาก AI)
type BestSeller struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
	Qty   int    `json:"qty"`
}

// SummaryResponse ตาม API Contract
// summary + recommendations มาจาก AI, best_sellers มาจากข้อมูลจริง
type SummaryResponse struct {
	Summary         string       `json:"summary"`
	BestSellers     []BestSeller `json:"best_sellers"`
	Recommendations []string     `json:"recommendations"`
}
