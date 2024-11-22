package dal

import (
	"context"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	redis2 "github.com/ksel172/Meduza/teamserver/internal/storage/repos"
)

type CheckInDAL struct {
	redis redis2.Service
}

func NewCheckInDAL(redisService *redis2.Service) *CheckInDAL {
	return &CheckInDAL{redis: *redisService}
}
func (dal *CheckInDAL) CreateAgent(agent models.Agent) error {
	if err := dal.redis.JsonSet(context.Background(), agent.RedisID(), agent); err != nil {
		return fmt.Errorf("Failed to register agent: %w", err)
	}
	return nil
}
