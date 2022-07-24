package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/repositories"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
	"github.com/zhanbolat18/parcel/users/pkg/crypto"
)

type ManageUser struct {
	hasher crypto.PasswordHasher
	repo   repositories.UserRepository
}

func NewUserService(hasher crypto.PasswordHasher, repo repositories.UserRepository) *ManageUser {
	return &ManageUser{hasher: hasher, repo: repo}
}

func (m *ManageUser) Get(ctx context.Context, id uint) (*entities.User, error) {
	return m.repo.GetById(ctx, id)
}

func (m *ManageUser) SignUp(ctx context.Context, email, password string) (*entities.User, error) {
	u, err := m.createUser(ctx, email, password, valueobjects.User)
	if err != nil {
		return nil, fmt.Errorf("create courier: %w", err)
	}
	return u, err
}

func (m *ManageUser) CreateCourier(ctx context.Context, email, password string) (*entities.User, error) {
	u, err := m.createUser(ctx, email, password, valueobjects.Courier)
	if err != nil {
		return nil, fmt.Errorf("create courier: %w", err)
	}
	return u, err
}

func (m *ManageUser) Couriers(ctx context.Context) ([]*entities.User, error) {
	users, err := m.repo.GetAllByRole(ctx, valueobjects.Courier)
	if err != nil {
		return nil, fmt.Errorf("get couriers: %w", err)
	}
	return users, nil
}

func (m *ManageUser) Courier(ctx context.Context, id uint) (*entities.User, error) {
	users, err := m.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get couriers: %w", err)
	}
	if users == nil || users.Role != valueobjects.Courier {
		return nil, errors.New(fmt.Sprintf("courier with id \"%d\" not found", id))
	}
	return users, nil
}

func (m *ManageUser) createUser(ctx context.Context, email, password string, role valueobjects.Role) (*entities.User, error) {
	u, err := m.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email \"%s\": %w", email, err)
	}
	if u != nil {
		return u, errors.New(fmt.Sprintf("user with email \"%s\" already exists", email))
	}
	pwdHash, err := m.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("create password hash: %w", err)
	}
	u = entities.NewUser(email, string(pwdHash), role)
	err = m.repo.Save(ctx, u)
	return u, err
}
