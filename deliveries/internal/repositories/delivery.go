package repositories

import (
	"context"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
)

type DeliveriesRepository interface {
	GetById(ctx context.Context, id uint) (*entities.Delivery, error)
	GetAll(ctx context.Context) ([]entities.Delivery, error)
	AssignToCourier(ctx context.Context, delivery *entities.Delivery, courier *entities.User) error
	Store(ctx context.Context, delivery *entities.Delivery) error
	Update(ctx context.Context, delivery *entities.Delivery) error
}
