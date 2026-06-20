package product

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

// shopID ดึง shop_id (hex) ที่ JWT middleware set ไว้ แล้วแปลงเป็น ObjectID
func shopID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.GetString(middleware.CtxShopID))
	if err != nil {
		return primitive.NilObjectID, false
	}
	return id, true
}

// List — GET /api/v1/products
func (h *Handler) List(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	data, err := h.service.List(c.Request.Context(), sid)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}

// Create — POST /api/v1/products
func (h *Handler) Create(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := h.service.Create(c.Request.Context(), sid, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}

// Update — PUT /api/v1/products/:id
func (h *Handler) Update(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := h.service.Update(c.Request.Context(), sid, c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(c)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}

// Delete — DELETE /api/v1/products/:id
func (h *Handler) Delete(c *gin.Context) {
	sid, ok := shopID(c)
	if !ok {
		response.Unauthorized(c)
		return
	}
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), sid, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(c)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, DeleteResponse{ID: id, Deleted: true})
}
