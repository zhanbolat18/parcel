package repositories

import (
	"context"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
)

type DeliveriesRepository interface {
	GetAll(ctx context.Context) ([]*entities.Delivery, error)
	GetAllByCourier(ctx context.Context, courierId uint) ([]*entities.Delivery, error)
	GetById(ctx context.Context, id uint) (*entities.Delivery, error)
	Store(ctx context.Context, delivery *entities.Delivery) error
	Update(ctx context.Context, delivery *entities.Delivery) error
}
