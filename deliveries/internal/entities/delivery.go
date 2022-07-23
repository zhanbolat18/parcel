package entities

import (
	"github.com/zhanbolat18/parcel/deliveries/internal/valueobjects"
	"time"
)

type Delivery struct {
	Id          uint                `json:"id,omitempty"`
	Status      valueobjects.Status `json:"status" json:"status,omitempty"`
	Destination string              `json:"destination" json:"destination,omitempty"`
	Recipient   *User               `json:"recipient" json:"recipient,omitempty"`
	Courier     *User               `json:"courier,omitempty" json:"courier,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

func NewDelivery(destination string, recipient *User) *Delivery {
	return &Delivery{
		Destination: destination,
		Recipient:   recipient,
		Status:      valueobjects.Created,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
