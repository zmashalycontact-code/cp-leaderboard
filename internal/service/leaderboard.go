package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/redis/go-redis/v9"
	"github.com/zmashaly/cp-leaderboard/internal/models"
)

// تأكد إن الـ Key هنا هو نفس اللي بتستخدمه في الـ Sync Engine
const leaderboardKey = "leaderboard"

type LeaderboardService struct {
	redis *redis.Client
}

func NewLeaderboardService(rdb *redis.Client) *LeaderboardService {
	return &LeaderboardService{redis: rdb}
}

// الدالة دي بقت أذكى، بتسحب كل اليوزرز كـ JSON وتفكهم وترتبهم
func (s *LeaderboardService) GetLeaderboard(ctx context.Context, limit int64) ([]models.User, error) {
	// 1. قراءة كل البيانات من الـ Hash Set
	// بنستخدم HGetAll لأننا بنخزن Handle -> UserJSON
	data, err := s.redis.HGetAll(ctx, leaderboardKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leaderboard from redis: %w", err)
	}

	users := make([]models.User, 0, len(data))
	for _, jsonStr := range data {
		var u models.User
		if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {

			continue
		}
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].SeasonPoints > users[j].SeasonPoints
	})

	if limit > 0 && int64(len(users)) > limit {
		users = users[:limit]
	}

	return users, nil
}
