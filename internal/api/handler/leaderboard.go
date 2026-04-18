package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmashaly/cp-leaderboard/internal/service"
)

type LeaderboardHandler struct {
	service *service.LeaderboardService
}

func NewLeaderboardHandler(s *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		service: s,
	}
}

// GET /api/leaderboard?limit=50
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	ctx := c.Request.Context()

	// --- Parse query param ---
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit parameter",
		})
		return
	}

	// --- Fetch from Redis ONLY ---
	entries, err := h.service.GetLeaderboard(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch leaderboard",
			"details": err.Error(),
		})
		return
	}

	// --- Success response ---
	c.JSON(http.StatusOK, gin.H{
		"data":  entries,
		"count": len(entries),
	})
}
