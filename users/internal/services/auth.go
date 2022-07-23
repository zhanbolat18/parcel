package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/repositories"
	. "github.com/zhanbolat18/parcel/users/internal/valueobjects"
	"github.com/zhanbolat18/parcel/users/pkg/crypto"
	"github.com/zhanbolat18/parcel/users/pkg/jwt"
)

type AuthService struct {
	hasher crypto.PasswordHasher
	jwt    jwt.Jwt
	repo   repositories.UserRepository
}

func NewAuthService(hasher crypto.PasswordHasher, jwt jwt.Jwt, repo repositories.UserRepository) *AuthService {
	return &AuthService{hasher: hasher, jwt: jwt, repo: repo}
}

func (a *AuthService) Authentication(ctx context.Context, email, password string) (token string, err error) {
	u, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("get by email \"%s\": %w", email, err)
	}
	if u == nil {
		return "", errors.New(fmt.Sprintf("user with email \"%s\" not found", email))
	}
	if !a.hasher.ComparePassword(password, u.PasswordHash) {
		return "", errors.New("invalid password")
	}
	if !a.validStatus(u) {
		return "", errors.New("access denied")
	}

	token, err = a.jwt.Generate(u.Id, u.Email)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return token, nil
}

func (a *AuthService) Authorization(ctx context.Context, token string) (*entities.User, error) {
	uId, _, err := a.jwt.Parse(token)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	u, err := a.repo.GetById(ctx, uId)
	if err != nil {
		return nil, fmt.Errorf("get user by id \"%d\": %w", uId, err)
	}
	return u, nil
}

func (a *AuthService) validStatus(user *entities.User) bool {

	switch user.Status {
	case Blocked, Frozen:
		return false
	case Active:
		return true
	default:
		return false
	}
}
