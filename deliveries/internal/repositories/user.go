package repositories

import (
	"context"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
)

type UsersRepository interface {
	GetCourier(ctx context.Context, id uint) (*entities.User, error)
	GetRecipient(ctx context.Context, id uint) (*entities.User, error)
}
