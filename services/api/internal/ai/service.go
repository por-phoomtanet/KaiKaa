package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kaikaa/api/internal/report"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	reports *report.Service
	llm     Completer
}

func NewService(reports *report.Service, llm Completer) *Service {
	return &Service{reports: reports, llm: llm}
}

// Summary ดึงยอดขายของวัน → ส่ง prompt ให้ AI → คืนสรุป + คำแนะนำ
// best_sellers มาจากข้อมูลจริง (ไม่พึ่ง AI) เพื่อความถูกต้อง
func (s *Service) Summary(ctx context.Context, shopID primitive.ObjectID, dateStr string) (SummaryResponse, error) {
	daily, err := s.reports.Daily(ctx, shopID, dateStr)
	if err != nil {
		return SummaryResponse{}, err // เช่น ErrInvalidDate → handler map 400
	}

	best := make([]BestSeller, 0, 3)
	for i, b := range daily.BestSellers {
		if i >= 3 {
			break
		}
		best = append(best, BestSeller{Name: b.Name, Emoji: b.Emoji, Qty: b.Qty})
	}

	raw, err := s.llm.Complete(ctx, buildPrompt(daily, best))
	if err != nil {
		return SummaryResponse{}, fmt.Errorf("ai upstream: %w", err)
	}

	parsed, err := parseAIJSON(raw)
	if err != nil {
		return SummaryResponse{}, fmt.Errorf("ai response invalid: %w", err)
	}

	return SummaryResponse{
		Summary:         parsed.Summary,
		BestSellers:     best,
		Recommendations: parsed.Recommendations,
	}, nil
}

// buildPrompt ประกอบ prompt ภาษาไทย สั่งให้ AI ตอบเป็น JSON เท่านั้น
func buildPrompt(d report.DailyReport, best []BestSeller) string {
	var sb strings.Builder
	sb.WriteString("คุณเป็นผู้ช่วยวิเคราะห์ยอดขายร้านกาแฟ ")
	sb.WriteString("วิเคราะห์ข้อมูลด้านล่างแล้วตอบกลับเป็น JSON ภาษาไทยเท่านั้น ")
	sb.WriteString(`รูปแบบ: {"summary":"สรุปยอดขาย 1-2 ประโยค","recommendations":["คำแนะนำข้อ1","ข้อ2","ข้อ3"]} `)
	sb.WriteString("ห้ามมีข้อความอื่นนอกเหนือจาก JSON\n\n")
	sb.WriteString(fmt.Sprintf("ยอดขายรวม: %d บาท จาก %d รายการ\n", d.Total, d.Count))
	sb.WriteString(fmt.Sprintf("เงินสด: %d บาท (%d%%), โอนเงิน: %d บาท (%d%%)\n", d.Cash, d.CashPct, d.Transfer, d.TransferPct))
	if len(best) > 0 {
		sb.WriteString("สินค้าขายดี: ")
		parts := make([]string, 0, len(best))
		for _, b := range best {
			parts = append(parts, fmt.Sprintf("%s %s (%d แก้ว)", b.Emoji, b.Name, b.Qty))
		}
		sb.WriteString(strings.Join(parts, ", ") + "\n")
	}
	return sb.String()
}

type aiPayload struct {
	Summary         string   `json:"summary"`
	Recommendations []string `json:"recommendations"`
}

// parseAIJSON ดึง JSON object ออกจากคำตอบ AI (เผื่อมี markdown fence หรือข้อความหุ้ม)
func parseAIJSON(raw string) (aiPayload, error) {
	s := strings.TrimSpace(raw)
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end < start {
		return aiPayload{}, fmt.Errorf("no JSON object found in AI response")
	}
	s = s[start : end+1]

	var p aiPayload
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return aiPayload{}, err
	}
	return p, nil
}
