package repo

import (
	"context"

	"github.com/crestenstclair/crud/internal/user"
)

//go:generate mockery --name Repo
type Repo interface {
	GetUser(ctx context.Context, userID string) (*user.User, error)
	//  DeleteUser(ctx context.Context, userID string) (User, error)
	// UpdateUser(ctx context.Context, user User) (User, error)
	// CreateUser(ctx context.Context, user User) (User, error)
}
