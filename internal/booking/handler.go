package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// POST /api/drafts
func (h *Handler) CreateDraft(c *gin.Context) {
	var req CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	id, err := h.svc.CreateDraft(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save draft"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "draft_id": id})
}

// GET /api/drafts/:id
func (h *Handler) GetDraft(c *gin.Context) {
	id := c.Param("id")
	draft, err := h.svc.GetDraft(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found or expired"})
		return
	}

	// 注意：直接回傳 draft 會包含所有欄位，如果只需要 draft_data 可以只回傳 draft.DraftData
	c.JSON(http.StatusOK, draft)
}
