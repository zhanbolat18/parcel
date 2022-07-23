package repositories

import "github.com/zhanbolat18/parcel/deliveries/internal/entities"

type UsersRepository interface {
	GetCourier(id uint) (*entities.User, error)
	GetRecipient(id uint) (*entities.User, error)
}
