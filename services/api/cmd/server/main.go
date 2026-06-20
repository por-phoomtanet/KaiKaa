package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/config"
	"github.com/kaikaa/api/internal/auth"
	"github.com/kaikaa/api/internal/middleware"
	"github.com/kaikaa/api/internal/response"
)

func main() {
	cfg := config.Load()
	db := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)

	// wire dependencies (handler ← service ← db)
	authService := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		// public — ไม่ต้อง token
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// protected — ต้องผ่าน JWT middleware (products/sales/... ต่อใน task ถัดไป)
		v1 := api.Group("/v1")
		v1.Use(middleware.JWT(cfg.JWTSecret))
		{
			// endpoint ตัวอย่างไว้ยืนยันว่า auth ทำงาน + ใช้ทดสอบ DoD (401 ถ้าไม่มี token)
			v1.GET("/me", func(c *gin.Context) {
				response.OK(c, gin.H{
					"user_id": c.GetString(middleware.CtxUserID),
					"shop_id": c.GetString(middleware.CtxShopID),
				})
			})
		}
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
