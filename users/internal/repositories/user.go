package repositories

import (
	"context"
	"github.com/zhanbolat18/parcel/users/internal/entities"
)

type UserRepository interface {
	GetById(ctx context.Context, id uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Save(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
}
