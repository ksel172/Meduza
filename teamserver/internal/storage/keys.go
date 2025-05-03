package storage

import "sync"

// The implementations in this package are just placeholders

// Singleton implementations for now
var (
	KeyRegistry = keyRegistry{
		registry: make(map[string][]byte),
	}
	AsymmetricKeyRegistry = asymmetricKeyRegistry{
		registry: make(map[string]KeyPair),
	}
)

// keyRegistry is used for mapping agent session tokens to AES keys
type keyRegistry struct {
	mu       sync.Mutex
	registry map[string][]byte
}

// asymmetricKeyStore is used for mapping agent auth token to ECDH asymmetric key
// for agent's first authentication and session token creation

// The Controller is responsible for writing all of the agents keys into the storage
// when the controller is added
type asymmetricKeyRegistry struct {
	mu       sync.Mutex
	registry map[string]KeyPair
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

func (k *asymmetricKeyRegistry) WriteKey(authToken string, keyPair KeyPair) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.registry[authToken] = keyPair
}

func (k *asymmetricKeyRegistry) GetKeys(authToken string) (KeyPair, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	key, exists := k.registry[authToken]
	return key, exists
}

func (k *asymmetricKeyRegistry) DeleteKey(authToken string) {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.registry, authToken)
}
