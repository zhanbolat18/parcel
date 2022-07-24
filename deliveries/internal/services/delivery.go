package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	"github.com/zhanbolat18/parcel/deliveries/internal/repositories"
	"github.com/zhanbolat18/parcel/deliveries/internal/valueobjects"
	"time"
)

type ManageDelivery struct {
	deliveryRepo repositories.DeliveriesRepository
	usersRepo    repositories.UsersRepository
}

func NewManageDelivery(deliveryRepo repositories.DeliveriesRepository, usersRepo repositories.UsersRepository) *ManageDelivery {
	return &ManageDelivery{deliveryRepo: deliveryRepo, usersRepo: usersRepo}
}

func (m *ManageDelivery) Create(ctx context.Context, recipient *entities.User, destination string) (*entities.Delivery, error) {
	delivery := entities.NewDelivery(destination, recipient)
	err := m.deliveryRepo.Store(ctx, delivery)
	if err != nil {
		return nil, fmt.Errorf("store delivery %w", err)
	}
	return delivery, nil
}

func (m *ManageDelivery) GetAll(ctx context.Context) ([]*entities.Delivery, error) {
	deliveries, err := m.deliveryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch all deliveries %w", err)
	}
	return deliveries, nil
}

func (m *ManageDelivery) GetOne(ctx context.Context, id uint) (*entities.Delivery, error) {
	delivery, err := m.deliveryRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get one delivery %w", err)
	}
	return delivery, nil
}

func (m *ManageDelivery) GetAllByCourier(ctx context.Context, courier *entities.User) ([]*entities.Delivery, error) {
	deliveries, err := m.deliveryRepo.GetAllByCourier(ctx, courier.Id)
	if err != nil {
		return nil, fmt.Errorf("fetch all deliveries %w", err)
	}
	return deliveries, nil
}

func (m *ManageDelivery) AssignToCourier(ctx context.Context, deliveryId, courierId uint) (*entities.Delivery, error) {
	courier, err := m.usersRepo.GetCourier(ctx, courierId)
	if err != nil {
		return nil, fmt.Errorf("get courier by id \"%d\": %w", courierId, err)
	}
	delivery, err := m.deliveryRepo.GetById(ctx, deliveryId)
	if err != nil {
		return nil, fmt.Errorf("get delivery by id \"%d\": %w", deliveryId, err)
	}
	if !m.isCourierAssignable(delivery) {
		return nil, errors.New("delivery is not assignable")
	}
	delivery.Status = valueobjects.Delivers
	delivery.UpdatedAt = time.Now()
	delivery.CourierId = &courier.Id
	err = m.deliveryRepo.Update(ctx, delivery)
	if err != nil {
		return nil, fmt.Errorf("assign delivery: %w", err)
	}
	return delivery, nil
}

func (m *ManageDelivery) Complete(ctx context.Context, deliveryId uint, courier *entities.User) (*entities.Delivery, error) {
	delivery, err := m.deliveryRepo.GetById(ctx, deliveryId)
	if err != nil {
		return nil, fmt.Errorf("get delivery by id \"%d\": %w", deliveryId, err)
	}
	if delivery.CourierId != nil && *delivery.CourierId != courier.Id {
		return nil, fmt.Errorf("forbidden")
	}
	if !m.isCompletable(delivery) {
		return nil, errors.New("delivery is not completable")
	}
	delivery.Status = valueobjects.Completed
	delivery.UpdatedAt = time.Now()
	err = m.deliveryRepo.Update(ctx, delivery)
	if err != nil {
		return nil, fmt.Errorf("update delivery: %w", err)
	}
	return delivery, nil
}

func (m *ManageDelivery) isCourierAssignable(delivery *entities.Delivery) bool {
	switch delivery.Status {
	case valueobjects.Created, valueobjects.Delivers:
		return true
	case valueobjects.Completed, valueobjects.Canceled:
		return false
	default:
		return false
	}
}

func (m *ManageDelivery) isCompletable(delivery *entities.Delivery) bool {
	return delivery.Status == valueobjects.Delivers
}
