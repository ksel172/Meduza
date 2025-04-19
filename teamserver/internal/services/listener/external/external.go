package external

import (
	"context"
)

// Create external listener type
// It should have:
// 1. An API to receives requests that have already been handled by some externally deployed listener
// 2. A client to send requests to the external listener
// 3. Tracking information of how to reach the listener, IP and Port

type ExternalListener struct {
	// host   string
	// port   int
	client *ExternalClient
}

func NewExternalListener(client *ExternalClient) (*ExternalListener, error) {
	return &ExternalListener{
		client: client,
	}, nil
}

func (l *ExternalListener) Start(ctx context.Context) error {
	return nil
}

func (l *ExternalListener) Stop(ctx context.Context) error {
	return nil
}

func (l *ExternalListener) Terminate(ctx context.Context) error {
	return nil
}

func (l *ExternalListener) UpdateConfig(ctx context.Context) error {
	return nil
}

func (l *ExternalListener) Validate() error {
	return nil
}
