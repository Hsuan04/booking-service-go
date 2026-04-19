package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	//dsn := "postgresql://postgres:haodjlc09121104@db.wulmfecstmeyegpqlevu.supabase.co:5432/postgres"

	// 取得在 Zeabur 設定的連線DB的環境變數
	dsn := os.Getenv("DATABASE_URL")
	fmt.Printf("嘗試連線的 DSN: %s\n", dsn)
	if dsn == "" {
		log.Fatal("DATABASE_URL 未設定，請檢查環境變數")
	}

	if strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		userInfo := parts[0]   // postgres:password
		serverInfo := parts[1] // host:port/db

		userParts := strings.Split(userInfo, "//")
		if len(userParts) > 1 {
			username := strings.Split(userParts[1], ":")[0]
			fmt.Printf("連線帳號 (User): %s\n", username)
		}
		fmt.Printf("目標主機 (Host): %s\n", serverInfo)
	}

	// 3. 網路探針：測試是否能連到該 Host 的 Port
	// 這能區分是「網路不通 (IPv6問題)」還是「帳密錯誤」
	testConnect(dsn)

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

func testConnect(dsn string) {
	// 簡單解析 host:port
	// 範例：aws-0-ap-northeast-1.pooler.supabase.com:5432
	// 這裡我們直接用最簡單的方式抓取
	fmt.Println("正在進行 TCP 握手測試...")

	// 這裡只是輔助，若解析失敗不影響主程式
	if strings.Contains(dsn, "@") && strings.Contains(dsn, ":") {
		hostPort := strings.Split(strings.Split(dsn, "@")[1], "/")[0]
		timeout := 5 * time.Second
		conn, err := net.DialTimeout("tcp", hostPort, timeout)
		if err != nil {
			fmt.Printf("網路警告：無法建立 TCP 連線到 %s, 錯誤: %v\n", hostPort, err)
			fmt.Println("提示：如果錯誤是 'network is unreachable'，代表是 IPv6/IPv4 打架。")
		} else {
			fmt.Printf("網路檢查：成功連線到 %s (TCP 通暢)\n", hostPort)
			_ = conn.Close()
		}
	}
}
