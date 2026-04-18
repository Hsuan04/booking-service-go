package main

import (
	"github.com/Hsuan04/booking-service-go/config" // 必須包含完整的 GitHub 路徑
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// 啟動資料庫
	config.ConnectDatabase()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "API 正常運作中，且已連線至 Supabase",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
