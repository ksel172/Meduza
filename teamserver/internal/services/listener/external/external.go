package external

import "sync"

// Create external listener type
// It should have:
// 1. An API to receives requests that have already been handled by some externally deployed listener
// 2. A client to send requests to the external listener
// 3. Tracking information of how to reach the listener, IP and Port

type ExternalListener struct {
	Config any

	isRunning bool
	mu        sync.RWMutex
}
