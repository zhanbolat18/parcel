package dto

type UserDto struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"min=8,required"`
}
