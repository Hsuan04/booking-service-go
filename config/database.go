package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 取得在 Zeabur 設定的連線DB的環境變數
	dsn := os.Getenv("DATABASE_URL")
	fmt.Printf("嘗試連線的 DSN: %s\n", dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL 未設定，請檢查環境變數")
	}

	// 配置日誌：這讓你可以在 Zeabur Logs 看到 Go 產生的每一條 SQL，對除錯極有幫助
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 執行超過一秒的 SQL 會被標記為警告
			LogLevel:      logger.Info, // 顯示所有 SQL
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("資料庫連線失敗: ", err)
	}

	// --- 高級連線池設定 ---
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("取得資料庫底層物件失敗: ", err)
	}

	// 設定最大閒置連線 (保持 10 個連線隨時待命，預約請求進來秒開)
	sqlDB.SetMaxIdleConns(10)

	// 設定最大開啟連線 (防止突然湧入的預約流量把 Supabase 方案的 Connection Limit 撐爆)
	sqlDB.SetMaxOpenConns(100)

	// 設定連線最長生命週期 (每小時更換一次連線，防止資料庫端無故斷線)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println(" Supabase 連線成功，連線池已啟動！")
	DB = db
}
