package redis

import (
	"context"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type CheckInDAL struct {
	redis Service
}

func NewCheckInDAL(redisService *Service) *CheckInDAL {
	return &CheckInDAL{redis: *redisService}
}
func (dal *CheckInDAL) CreateAgent(agent models.Agent) error {
	if _, err := dal.redis.JsonSet(context.Background(), agent.ID, agent); err != nil {
		return fmt.Errorf("Failed to register agent: %w", err)
	}
	return nil
}
