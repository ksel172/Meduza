package dal

import (
	"context"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
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
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)

	if err := pubsub.Ping(ctx, ""); err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel: %v", err)
	}
	messageChan := make(chan string)

	go func() {
		defer pubsub.Close()
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

func (dal *PubSubDAL) Publish(ctx context.Context, channel string, message string) error {
	rdb := dal.client.GetClient()

	listKey := fmt.Sprintf("chat_history:%s", channel)
	if err := rdb.RPush(ctx, listKey, message).Err(); err != nil {
		return fmt.Errorf("failed to store message in history: %v", err)
	}

	if err := rdb.Publish(ctx, channel, message).Err(); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (dal *PubSubDAL) ReceiveMessage(ctx context.Context, channel string) (string, error) {
	rdb := dal.client.GetClient()
	pubsub := rdb.Subscribe(ctx, channel)

	if err := pubsub.Ping(ctx, ""); err != nil {
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

	listKey := fmt.Sprintf("chat_history:%s", channel)
	messages, err := rdb.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chat history: %v", err)
	}

	return messages, nil
}
