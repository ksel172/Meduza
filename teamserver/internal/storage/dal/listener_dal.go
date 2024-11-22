package dal

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	redis "github.com/ksel172/Meduza/teamserver/internal/storage/repos"
)

const listenerPrefix = "listener:"

type ListenerDAL struct {
	redis redis.Service
}

func NewListenerDAL(redis *redis.Service) *ListenerDAL {
	return &ListenerDAL{*redis}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	if listener.ID == "" {
		listener.ID = uuid.New().String()
	}

	key := fmt.Sprintf("%s%s", listenerPrefix, listener.ID)
	data, err := json.MarshalContext(ctx, listener)
	if err != nil {
		return fmt.Errorf("failed to marshal listener: %w", err)
	}

	if err := dal.redis.JsonSet(ctx, key, data); err != nil {
		return fmt.Errorf("failed to set listener: %w", err)
	}
	return nil
}

func (dal *ListenerDAL) GetListener(ctx context.Context, id string) (*models.Listener, error) {
	key := fmt.Sprintf("%s%s", listenerPrefix, id)

	data, err := dal.redis.JsonGet(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get listener: %w", err)
	}

	var listener models.Listener
	if err := json.UnmarshalContext(ctx, []byte(data), &listener); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listener: %w", err)
	}
	return &listener, nil
}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, listener *models.Listener) error {
	return dal.CreateListener(ctx, listener)
}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, id string) error {
	key := fmt.Sprintf("%s%s", listenerPrefix, id)
	return dal.redis.JsonDelete(ctx, key)
}
