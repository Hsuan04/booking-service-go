package booking

import (
	"encoding/json"
	"time"
)

// BookingDraft 實體
type BookingDraft struct {
	ID        string          `json:"id"`
	StudioID  string          `json:"studio_id"`
	DraftData json.RawMessage `json:"draft_data"` // 存儲原始 JSON 內容
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at"`
}

// CreateDraftRequest DTO
type CreateDraftRequest struct {
	StudioID  string          `json:"studio_id" binding:"required"`
	DraftData json.RawMessage `json:"draft_data" binding:"required"`
}
