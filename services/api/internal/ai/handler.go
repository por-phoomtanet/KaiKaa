package ai

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/internal/middleware"
	"github.com/kaikaa/api/internal/report"
	"github.com/kaikaa/api/internal/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func shopID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.GetString(middleware.CtxShopID))
	if err != nil {
		return primitive.NilObjectID, false
	}
	return id, true
}

// Summary — POST /api/v1/ai/summary
func (h *Handler) Summary(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}

	// body เป็น optional ({} หรือ {date}) — ยอมรับ body ว่าง (EOF)
	var req SummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.service.Summary(c.Request.Context(), sid, req.Date)
	if err != nil {
		if errors.Is(err, report.ErrInvalidDate) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		// AI ล่ม/timeout/parse ไม่ได้ → 502 (อ่านได้ ไม่ crash)
		response.Error(c, http.StatusBadGateway, err.Error())
		return
	}
	response.OK(c, data)
}
