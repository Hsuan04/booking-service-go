package database

import (
	"database/sql"
	"fmt"
	"time"

	_ ".github.com/lib/pq" // 引入 PostgreSQL 驅動 (底線表示只觸發 init()，不直接呼叫函數)
)

// Config 定義了資料庫連線所需的參數，這使得它具備了「多環境彈性」
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB 負責建立並回傳一個設定好連線池的 DB 實例
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	// 組裝 DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// sql.Open 只是驗證參數，不會立刻建立真實連線
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("無法開啟資料庫連線: %w", err)
	}

	// db.Ping 才會真正去連線資料庫，確保設定是正確的
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("資料庫 ping 失敗: %w", err)
	}

	// === 設定專業的連線池 (Connection Pool) ===

	// 1. 最大開啟連線數：避免突發流量把資料庫打掛 (依據你的 DB 規格調整)
	db.SetMaxOpenConns(25)

	// 2. 最大閒置連線數：保留一定數量的連線以應付快速發生的請求，通常與 MaxOpenConns 相同或略小
	db.SetMaxIdleConns(25)

	// 3. 連線最長生命週期：定期強制關閉並重建連線，避免遇到網路設備(如 Firewall/Load Balancer)強行切斷閒置連線
	db.SetConnMaxLifetime(15 * time.Minute)

	// 4. 閒置連線的最長生命週期：若連線一直沒人用，多久後關閉回收資源
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
