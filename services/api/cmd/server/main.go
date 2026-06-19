package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/config"
)

func main() {
	cfg := config.Load()

	// ต่อ MongoDB (ping ในตัว) — ใช้ db ใน task ถัดไป (auth/product/...)
	db := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)
	_ = db

	r := gin.Default()

	// health check — ใช้ตรวจ DoD ของ T01
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
