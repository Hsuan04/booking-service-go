package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 路由器
	r := gin.Default()

	// 建立一個簡單的 GET API 端點
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Zeabur! Lawrence 的預約微服務 API 啟動成功！",
			"status":  "success",
		})
	})

	// 獲取系統環境變數的 PORT (Zeabur 部署時需要)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 本機開發時預設使用 8080
	}

	log.Printf("伺服器準備在 Port %s 啟動...\n", port)

	// 啟動伺服器
	r.Run(":" + port)
}
