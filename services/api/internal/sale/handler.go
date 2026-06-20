package sale

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/internal/middleware"
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

// Create — POST /api/v1/sales
func (h *Handler) Create(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	var req CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// รวมถึง method ที่ไม่ใช่ cash/transfer (binding oneof) → 400
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := h.service.Create(c.Request.Context(), sid, req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}

// List — GET /api/v1/sales?date=YYYY-MM-DD
func (h *Handler) List(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	data, err := h.service.ListByDate(c.Request.Context(), sid, c.Query("date"))
	if err != nil {
		if errors.Is(err, ErrInvalidDate) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}
