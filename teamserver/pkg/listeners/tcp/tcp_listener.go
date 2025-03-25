package tcp_listener

import (
	"context"
)

type TCPListener struct {
}

func (l TCPListener) Start(ctx context.Context) error {
	return nil
}
func (l TCPListener) Stop(ctx context.Context) error {
	return nil
}
func (l TCPListener) Terminate(ctx context.Context) error {
	return nil
}
func (l TCPListener) UpdateConfig(ctx context.Context) error {
	return nil
}
