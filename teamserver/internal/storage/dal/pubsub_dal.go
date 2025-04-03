package dal

import (
	"context"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type IPubSubDAL interface {
	Subscribe(ctx context.Context, channel string) (chan string, error)
	Unsubscribe(ctx context.Context, channel string) error
	Publish(ctx context.Context, channel string, message string) error
}

type PubSubDAL struct {
	client repos.Service
}

func NewPubSubDAL(client repos.Service) *PubSubDAL {
	return &PubSubDAL{
		client: client,
	}
}

func (dal *PubSubDAL) Subscribe(ctx context.Context, channel string) (chan string, error) {
	logger.Info("INFO", "YOU GOT THIS FAR")
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)
	if err := pubsub; err != nil {
		logger.Info("ERROR", err)
		return nil, fmt.Errorf("failed to subscribe to channel: %v", err)
	}

	messageChan := make(chan string)

	go func() {
		for msg := range pubsub.Channel() {
			messageChan <- msg.Payload
		}
	}()

	return messageChan, nil
}

func (dal *PubSubDAL) Unsubscribe(ctx context.Context, channel string) error {
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)
	if err := pubsub.Unsubscribe(ctx, channel); err != nil {
		return fmt.Errorf("failed to unsubscribe from channel: %v", err)
	}
	return nil
}

func (dal *PubSubDAL) PublishMessage(ctx context.Context, channel string, message string) error {
	rdb := dal.client.GetClient()
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (dal *PubSubDAL) ReceiveMessage(ctx context.Context, channel string) (string, error) {
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)
	if err := pubsub; err != nil {
		return "", fmt.Errorf("failed to subscribe to channel: %v", err)
	}

	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to receive message: %w", err)
	}

	return msg.Payload, nil
}

func (dal *PubSubDAL) GetAllMessages(ctx context.Context, channel string) ([]string, error) {
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)
	if err := pubsub; err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel: %v", err)
	}

	messages := []string{}
	for msg := range pubsub.Channel() {
		messages = append(messages, msg.Payload)
	}

	return messages, nil
}
