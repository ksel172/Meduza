package services

import (

	//standard
	"context"
	"database/sql"
	"errors"
)

type IListenerDAL interface {
	CreateListener(context.Context, *Listener) error
	GetListenerById(context.Context, string) (Listener, error)
	GetAllListeners(context.Context) ([]Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
	GetActiveListeners(context.Context) ([]Listener, error)
	GetListenerByName(context.Context, string) (Listener, error)
}

var (
	logLevel          = "[DAL]"
	logDetailListener = "[Listener]"
	ErrUnimplemented  = errors.New("method not implemented")
)

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

// TODO: Also have to change the .sql for these

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *Listener) error {
	return ErrUnimplemented
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (Listener, error) {
	return Listener{}, ErrUnimplemented
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]Listener, error) {
	return nil, ErrUnimplemented
}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, listenerID string) error {
	return ErrUnimplemented
}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, listenerID string, updates map[string]any) error {
	return ErrUnimplemented
}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]Listener, error) {
	return nil, ErrUnimplemented
}

func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (Listener, error) {
	return Listener{}, ErrUnimplemented
}
