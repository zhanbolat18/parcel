package repositories

import (
	"context"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
)

type UserRepository interface {
	GetById(ctx context.Context, id uint) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetAllByRole(ctx context.Context, role valueobjects.Role) ([]*entities.User, error)
	Save(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
}
