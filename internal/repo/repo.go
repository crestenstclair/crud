package repo

import (
	"context"
)


type User struct {

}

//go:generate mockery --name Repo
type Repo interface {
  GetUser(ctx context.Context, userID string) (User, error)
  DeleteUser(ctx context.Context, userID string) (User, error)
  UpdateUser(ctx context.Context, user User) (User, error)
  CreateUser(ctx context.Context, user User) (User, error)
}
