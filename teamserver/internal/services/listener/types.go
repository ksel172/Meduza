package listener

import (
	"context"
	"sync"

	_ "github.com/go-playground/validator/v10"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type statusUpdate struct {
	listenerID string
	status     string
}

// ListenerService is the entrypoint for listener operations exposed to clients
type ListenerService struct {
	startTimeout int
	stopTimeout  int

	checkinController checkin.ICheckInController // not used directly, but injected into listener implementations
	listenerDal       dal.IListenerDAL

	// Keep track of the runtime listener representations
	activeListeners map[string]*Listener
	statusUpdates   chan statusUpdate
	mux             sync.RWMutex

	rootCtx context.Context
}

// Listener is a representation of a listener of any kind
type Listener struct {
	models.Listener // embed all of the data fields in the listener data model
	listener        ListenerImplementation
	mux             sync.RWMutex
	statusUpdatesCh chan string
}

func NewListenerService(listenerDAL dal.IListenerDAL, checkinController checkin.ICheckInController) *ListenerService {
	ls := &ListenerService{
		checkinController: checkinController,
		startTimeout:      30,
		stopTimeout:       30,
		listenerDal:       listenerDAL,
		activeListeners:   make(map[string]*Listener),
		statusUpdates:     make(chan statusUpdate, 100),
		rootCtx:           context.Background(),
	}

	// Start the status update processor
	go ls.processStatusUpdates()

	return ls
}
