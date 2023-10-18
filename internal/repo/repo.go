package repo

import (
	"context"

	"github.com/crestenstclair/crud/internal/user"
)

//go:generate mockery --name Repo
type Repo interface {
	GetUser(context.Context, string) (*user.User, error)
	//  DeleteUser(ctx context.Context, userID string) (User, error)
	UpdateUser(context.Context, user.User) (*user.User, error)
	CreateUser(context.Context, user.User) (*user.User, error)
}
