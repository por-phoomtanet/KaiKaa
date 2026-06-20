package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// รูปแบบ response กลางตาม T00a — ทุก handler ต้องผ่าน helper เหล่านี้ ห้าม c.JSON ตรง

// OK ส่ง success พร้อม data
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// Error ส่ง fail พร้อม message + status code ที่กำหนด
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"success": false, "message": message})
}

// Unauthorized ใช้ใน middleware — abort chain ทันที
func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
}

// NotFound สำหรับ resource ที่ไม่เจอ
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
}
