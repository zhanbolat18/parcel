package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/services"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
	"github.com/zhanbolat18/parcel/users/pkg/crypto"
	jwt2 "github.com/zhanbolat18/parcel/users/pkg/jwt"
	"strconv"
	"testing"
	"time"
)

type mockUserRepo struct {
	memory map[string]*entities.User
}

func (m *mockUserRepo) GetById(ctx context.Context, id uint) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	return m.memory[email], nil
}

func (m *mockUserRepo) Save(ctx context.Context, user *entities.User) error {
	m.memory[user.Email] = user
	return nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *entities.User) error {
	m.memory[user.Email] = user
	return nil
}

var hasher = crypto.NewPasswordHasher(13)
var ctx = context.Background()
var jwt = jwt2.NewJwtManager(10*time.Second, 0, []byte("customKey"))

func TestManageUser_CreateUser(t *testing.T) {
	assrt := assert.New(t)
	repo := &mockUserRepo{memory: make(map[string]*entities.User)}
	email := "custom@email.com"
	password := "custompassword"
	srv := services.NewUserService(hasher, repo)
	u, err := srv.SignUp(ctx, email, password)
	assrt.Nil(err)
	assrt.NotEqual(u.PasswordHash, password)
	assrt.Equal(u.Role, valueobjects.User)
	assrt.True(hasher.ComparePassword(password, u.PasswordHash))
}

func TestManageUser_CreateCourier(t *testing.T) {
	assrt := assert.New(t)
	repo := &mockUserRepo{memory: make(map[string]*entities.User)}
	email := "custom@email.com"
	password := "custompassword"
	srv := services.NewUserService(hasher, repo)
	u, err := srv.CreateCourier(ctx, email, password)
	assrt.Nil(err)
	assrt.NotEqual(u.PasswordHash, password)
	assrt.Equal(u.Role, valueobjects.Courier)
	assrt.True(hasher.ComparePassword(password, u.PasswordHash))
}

func TestAuthService_AuthenticationSuccess(t *testing.T) {
	assrt := assert.New(t)
	email := "active@email.com"
	password := "custompassword"
	hash, _ := hasher.Hash(password)
	repo := &mockUserRepo{memory: map[string]*entities.User{
		email: {
			Id:           1,
			Email:        email,
			PasswordHash: string(hash),
			Status:       valueobjects.Active,
			Role:         valueobjects.User,
		},
	}}
	srv := services.NewAuthService(hasher, jwt, repo)
	token, err := srv.Authentication(ctx, email, password)
	assrt.Nil(err)
	assrt.NotEmpty(token)
	assrt.True(jwt.Validate(token))
}

func TestAuthService_AuthenticationFail(t *testing.T) {
	assrt := assert.New(t)
	password := "custompassword"
	hash, _ := hasher.Hash(password)
	repo := &mockUserRepo{memory: map[string]*entities.User{
		"frozen@email.com": {
			Id:           1,
			Email:        "frozen@email.com",
			PasswordHash: string(hash),
			Status:       valueobjects.Frozen,
			Role:         valueobjects.User,
		},
		"blocked@custom.com": {
			Id:           1,
			Email:        "blocked@custom.com",
			PasswordHash: string(hash),
			Status:       valueobjects.Blocked,
			Role:         valueobjects.User,
		},
	}}

	failCases := []struct {
		email, password string
	}{
		{
			email:    "",
			password: password,
		},
		{
			email:    "frozen@email.com",
			password: password,
		},
		{
			email:    "frozen@email.com", // because status is frozen
			password: password,
		},
		{
			email:    "blocked@custom.com", // because status is blocked
			password: password,
		},
	}

	srv := services.NewAuthService(hasher, jwt, repo)

	for i, failCase := range failCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			token, err := srv.Authentication(ctx, failCase.email, failCase.password)
			assrt.NotNil(err)
			assrt.False(jwt.Validate(token))
		})
	}

}
