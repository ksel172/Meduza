package storage

import "sync"

// The implementations in this package are just placeholders

// Singleton implementations for now
var (
	KeyRegistry = keyRegistry{
		registry: make(map[string][]byte),
	}
)

// keyRegistry is used for mapping agent session tokens to AES keys
type keyRegistry struct {
	mu       sync.Mutex
	registry map[string][]byte
}

type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

func (k *keyRegistry) WriteKey(sessionToken string, key []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.registry[sessionToken] = key
}

func (k *keyRegistry) GetKey(sessionToken string) ([]byte, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	key, exists := k.registry[sessionToken]
	return key, exists
}
