package report

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

// Daily — GET /api/v1/reports/daily?date=YYYY-MM-DD
func (h *Handler) Daily(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	data, err := h.service.Daily(c.Request.Context(), sid, c.Query("date"))
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

// Monthly — GET /api/v1/reports/monthly?month=YYYY-MM
func (h *Handler) Monthly(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	data, err := h.service.Monthly(c.Request.Context(), sid, c.Query("month"))
	if err != nil {
		if errors.Is(err, ErrInvalidMonth) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}
