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

func (l *ExternalListener) Start(context.Context) error {
	// if err := l.client.start,

	return nil
}

func (l *ExternalListener) Stop(context.Context) error {
	return nil
}

func (l *ExternalListener) Terminate(context.Context) error {
	return nil
}

func (l *ExternalListener) UpdateConfig(context.Context) error {
	return nil
}
