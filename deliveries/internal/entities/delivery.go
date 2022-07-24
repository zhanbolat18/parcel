package entities

import (
	"github.com/zhanbolat18/parcel/deliveries/internal/valueobjects"
	"time"
)

type Delivery struct {
	Id          uint                `json:"id"`
	Status      valueobjects.Status `json:"status"`
	Destination string              `json:"destination"`
	RecipientId uint                `json:"recipient_id"`
	CourierId   *uint               `json:"courier_id,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

func NewDelivery(destination string, recipient *User) *Delivery {
	return &Delivery{
		Destination: destination,
		RecipientId: recipient.Id,
		Status:      valueobjects.Created,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
