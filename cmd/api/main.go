// cmd/api/main.go
package main

import (
	"log"
	"os"
	"strings"

	// 必須與 go.mod 的 module 名稱開頭一致
	"booking-service-go/internal/booking"
	"booking-service-go/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 環境變數加載邏輯
	appEnv := getEnvOrDefault("APP_ENV", "local")
	if appEnv == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Println("提示：未找到 .env.local，嘗試載入標準 .env")
			_ = godotenv.Load() // 忽略錯誤，若無則使用系統變數
		}
	}

	// 2. 初始化資料庫
	dbConfig := database.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "secret"),
		DBName:   getEnvOrDefault("DB_NAME", "booking_db"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("無法啟動資料庫: %v", err)
	}
	defer db.Close()

	log.Printf("[%s] 資料庫連線成功！", appEnv)

	// 3. 依賴注入
	repo := booking.NewRepository(db)
	svc := booking.NewService(repo)
	handler := booking.NewHandler(svc)

	// 4. 設定 Gin 伺服器
	if appEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// --- 5. 核心：動態跨域處理 (CORS) ---
	// 從環境變數讀取允許的網域，例如：
	// ALLOW_ORIGINS=http://localhost:3000,https://booking-dev.vercel.app
	allowedOriginsStr := getEnvOrDefault("ALLOW_ORIGINS", "http://localhost:3000")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 檢查請求網域是否在白名單中
		isAllowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 6. 路由定義
	api := r.Group("/api")
	{
		api.POST("/drafts", handler.CreateDraft)
		api.GET("/drafts/:id", handler.GetDraft)
	}

	// 7. 啟動伺服器
	port := getEnvOrDefault("PORT", "8080")
	log.Printf("伺服器運行環境 [%s], 埠號 %s...", appEnv, port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
