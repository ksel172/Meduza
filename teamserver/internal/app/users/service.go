package users

import (
	"context"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
)

type Service struct {
	dal *dal.UserDAL
}

func NewService(dal *dal.UserDAL) *Service {
	return &Service{dal: dal}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	// Business logic (e.g., filtering, sorting) can be applied here
	return s.dal.GetUsers(ctx)
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	// Validation or preprocessing logic
	if user.Username == "" {
		return fmt.Errorf("user name cannot be empty")
	}

	return s.dal.CreateUser(ctx, user)
}
