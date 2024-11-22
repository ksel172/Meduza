package redis

import (
	"context"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/internal/models"
)

type CheckInDAL struct {
	redis Service
}

func NewCheckInDAL(redisService *Service) *CheckInDAL {
	return &CheckInDAL{redis: *redisService}
}
func (dal *CheckInDAL) CreateAgent(agent models.Agent) error {
	if err := dal.redis.JsonSet(context.Background(), agent.RedisID(), agent); err != nil {
		return fmt.Errorf("Failed to register agent: %w", err)
	}
	return nil
}
