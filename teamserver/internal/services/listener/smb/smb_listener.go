package smb_listener

import "context"

type SMBListener struct {
}

type SMBListenerConfig struct {
}

func (l SMBListener) Start(ctx context.Context) error {
	return nil
}
func (l SMBListener) Stop(ctx context.Context) error {
	return nil
}
func (l SMBListener) Terminate(ctx context.Context) error {
	return nil
}
func (l SMBListener) UpdateConfig(ctx context.Context) error {
	return nil
}
func (l SMBListener) Validate() error {
	return nil
}
