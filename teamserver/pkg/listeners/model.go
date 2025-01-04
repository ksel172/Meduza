package listeners

import "time"

type ListenerController interface {
	Start() error
	Stop(time.Duration) error
	GetName() string
}
