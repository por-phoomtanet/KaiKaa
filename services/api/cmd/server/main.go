package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaikaa/api/config"
	"github.com/kaikaa/api/internal/auth"
	"github.com/kaikaa/api/internal/middleware"
	"github.com/kaikaa/api/internal/product"
	"github.com/kaikaa/api/internal/report"
	"github.com/kaikaa/api/internal/response"
	"github.com/kaikaa/api/internal/sale"
)

func main() {
	cfg := config.Load()
	db := config.ConnectMongo(cfg.MongoURI, cfg.MongoDB)

	// wire dependencies (handler ← service ← db)
	authHandler := auth.NewHandler(auth.NewService(db, cfg.JWTSecret))
	productHandler := product.NewHandler(product.NewService(db))
	saleHandler := sale.NewHandler(sale.NewService(db))
	reportHandler := report.NewHandler(report.NewService(db))

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

			// products (T03)
			v1.GET("/products", productHandler.List)
			v1.POST("/products", productHandler.Create)
			v1.PUT("/products/:id", productHandler.Update)
			v1.DELETE("/products/:id", productHandler.Delete)

			// sales (T04)
			v1.POST("/sales", saleHandler.Create)
			v1.GET("/sales", saleHandler.List)

			// reports (T05)
			v1.GET("/reports/daily", reportHandler.Daily)
			v1.GET("/reports/monthly", reportHandler.Monthly)
		}
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
