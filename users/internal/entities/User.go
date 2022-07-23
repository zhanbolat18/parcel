package entities

import "github.com/zhanbolat18/parcel/users/internal/valueobjects"

type User struct {
	Id           uint                `json:"id"`
	Email        string              `json:"email"`
	PasswordHash string              `json:"-"`
	Status       valueobjects.Status `json:"status"`
	Role         valueobjects.Role   `json:"role"`
}

func NewUser(email, passwordHash string, role valueobjects.Role) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Status:       valueobjects.Active,
		Role:         role}
}
