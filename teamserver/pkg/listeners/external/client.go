package external_listener

import (
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

// TODO: External listener client for sending commands to external listener implementations

type ListenerClient struct {
	// TODO: Let's think about how to authenticate listeners
	apiKey string

	listenerDal dal.IListenerDAL
}

func NewListenerClient(apiKey string, listenerDal dal.IListenerDAL) *ListenerClient {
	return &ListenerClient{
		apiKey:      apiKey,
		listenerDal: listenerDal,
	}
}

// AddListener contacts the external listener controller to add a new listener
func (lc *ListenerClient) AddListener(listener models.Listener, url string) error {

	return nil
}

// AddListener contacts the external listener controller to update a listener
func (lc *ListenerClient) UpdateListener(listenerID string) error {

	return nil
}

// AddListener contacts the external listener controller to start a listener
func (lc *ListenerClient) StartListener(listenerID string) error {

	return nil
}

// AddListener contacts the external listener controller to stop a listener
func (lc *ListenerClient) StopListener(listenerID string) error {

	return nil
}

// AddListener contacts the external listener controller to terminate listener
func (lc *ListenerClient) TerminateListener(listenerID string) error {

	return nil
}
