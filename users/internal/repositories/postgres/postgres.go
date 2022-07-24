package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/repositories"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
)

type user struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repositories.UserRepository {
	return &user{db: db}
}

type userModel struct {
	Id           uint   `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
	Status       string `db:"status"`
}

func (u *user) GetById(ctx context.Context, id uint) (*entities.User, error) {
	um := &userModel{}
	err := u.db.GetContext(ctx, um, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u.hydrateToEntity(um), nil
}

func (u *user) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	um := &userModel{}
	err := u.db.GetContext(ctx, um, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u.hydrateToEntity(um), nil
}

func (u *user) GetAllByRole(ctx context.Context, role valueobjects.Role) ([]*entities.User, error) {
	um := make([]userModel, 0)
	err := u.db.SelectContext(ctx, um, "SELECT * FROM users WHERE role=$1", role)
	if err != nil {
		return nil, err
	}
	users := make([]*entities.User, 0, len(um))
	for _, model := range um {
		users = append(users, u.hydrateToEntity(&model))
	}
	return users, nil
}

func (u *user) Save(ctx context.Context, user *entities.User) error {
	var id int
	q := "INSERT INTO users(email, password_hash, role, status) VALUES($1, $2, $3, $4) RETURNING id"
	err := u.db.QueryRowContext(ctx, q, user.Email, user.PasswordHash, user.Role, user.Status).Scan(&id)
	if err != nil {
		return err
	}
	user.Id = uint(id)
	return nil
}

func (u *user) Update(ctx context.Context, user *entities.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *user) hydrateToEntity(user *userModel) *entities.User {
	return &entities.User{
		Id:           user.Id,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Status:       valueobjects.Status(user.Status),
		Role:         valueobjects.Role(user.Role),
	}
}
