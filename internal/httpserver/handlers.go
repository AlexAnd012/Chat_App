package httpserver

import (
	"chatapp/internal/data"
	"chatapp/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	repo         *storage.RedisRepo
	historyLimit int
}

func NewHandlers(repo *storage.RedisRepo, historyLimit int) *Handlers {
	return &Handlers{repo: repo, historyLimit: historyLimit}
}

type createMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// Подготавливаем строку
func clear(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func (h *Handlers) PostMessage(c *gin.Context) {
	var in createMessage
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	in.Username = clear(in.Username)
	in.Text = clear(in.Text)
	if l := len(in.Username); l < 2 || l > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be 2..32 chars"})
		return
	}
	if l := len(in.Text); l < 1 || l > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text must be 1..2000 chars"})
		return
	}

	msg := data.Message{
		Id:       uuid.NewString(),
		Username: in.Username,
		Text:     in.Text,
		Sendtime: time.Now().UTC(),
	}

	b, _ := json.Marshal(msg)

	ctx := c.Request.Context()
	if err := h.repo.AppendMessage(ctx, b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.TrimHistory(ctx, h.historyLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (h *Handlers) GetMessages(c *gin.Context) {
	limit := h.historyLimit
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= h.historyLimit {
			limit = n
		}
	}
	items, err := h.repo.RecentMessages(c, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// В Redis новые в начале, разворачиваем наоборот
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	c.Data(http.StatusOK, "application/json", []byte("["+strings.Join(items, ",")+"]"))
}
