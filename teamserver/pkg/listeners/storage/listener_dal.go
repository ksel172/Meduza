package storage

import (
	//standard
	"context"
	"database/sql"

	// internal
	"github.com/ksel172/Meduza/teamserver/models"
)

type IListenerDAL interface {
	CreateListener(context.Context, *models.Listener) error
	GetListenerById(context.Context, string) (models.Listener, error)
	GetAllListeners(context.Context) ([]models.Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
	GetActiveListeners(context.Context) ([]models.Listener, error)
	GetListenerByName(context.Context, string) (models.Listener, error)
}

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	// TODO: Implement this method
	return nil
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (models.Listener, error) {
	// TODO: Implement this method
	return models.Listener{}, nil
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]models.Listener, error) {
	// TODO: Implement this method
	return []models.Listener{}, nil
}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, listenerID string) error {
	// TODO: Implement this method
	return nil
}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, listenerID string, updates map[string]any) error {
	// TODO: Implement this method
	return nil
}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]models.Listener, error) {
	// TODO: Implement this method
	return []models.Listener{}, nil
}

func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (models.Listener, error) {
	// TODO: Implement this method
	return models.Listener{}, nil
}
