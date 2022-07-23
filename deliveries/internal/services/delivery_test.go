package services_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	"github.com/zhanbolat18/parcel/deliveries/internal/services"
	"github.com/zhanbolat18/parcel/deliveries/internal/valueobjects"
	"testing"
	"time"
)

type mockRepos struct {
	deliveries map[uint]*entities.Delivery
	users      map[uint]*entities.User
}

func (m *mockRepos) GetCourier(id uint) (*entities.User, error) {
	if u, ok := m.users[id]; ok {
		if u.Role == "courier" {
			return u, nil
		}
		return nil, errors.New("user is not courier")
	}
	return nil, errors.New("courier not found")
}

func (m *mockRepos) GetRecipient(id uint) (*entities.User, error) {
	if u, ok := m.users[id]; ok {
		if u.Role == "user" {
			return u, nil
		}
		return nil, errors.New("user is not courier")
	}
	return nil, errors.New("courier not found")
}

func (m *mockRepos) GetById(_ context.Context, id uint) (*entities.Delivery, error) {
	return m.deliveries[id], nil
}

func (m *mockRepos) GetAll(_ context.Context) ([]entities.Delivery, error) {
	return nil, nil
}

func (m *mockRepos) AssignToCourier(_ context.Context, d *entities.Delivery, u *entities.User) error {
	if dd, ok := m.deliveries[d.Id]; ok {
		dd.Courier = u
	}
	return nil
}

func (m *mockRepos) Store(_ context.Context, _ *entities.Delivery) error { return nil }

func (m *mockRepos) Update(_ context.Context, _ *entities.Delivery) error { return nil }

var ctx = context.Background()

func TestManageDelivery_Create(t *testing.T) {
	repo := &mockRepos{}
	srv := services.NewManageDelivery(repo, repo)
	asrt := assert.New(t)
	recip := &entities.User{
		Id:    1,
		Email: "email@mail.com",
		Role:  "user",
	}
	d, err := srv.Create(ctx, recip, "Some Address 1, 14")
	asrt.Nil(err)
	asrt.NotNil(d)
	asrt.Equal(d.Recipient.Id, recip.Id)
	asrt.Nil(d.Courier)
	asrt.LessOrEqual(time.Now().Sub(d.CreatedAt).Milliseconds(), int64(1))
	asrt.LessOrEqual(time.Now().Sub(d.UpdatedAt).Milliseconds(), int64(1))
}

func TestManageDelivery_AssignToCourier(t *testing.T) {
	users := map[uint]*entities.User{
		1: {Id: 1, Email: "custom1@mail.com", Role: "courier"},
		4: {Id: 4, Email: "custom4@mail.com", Role: "user"},
	}
	deliveries := map[uint]*entities.Delivery{
		1: {Id: 1, Status: valueobjects.Created, Recipient: users[1]},
		2: {Id: 2, Status: valueobjects.Canceled, Recipient: users[1]},
		3: {Id: 3, Status: valueobjects.Delivers, Recipient: users[1]},
		4: {Id: 4, Status: valueobjects.Completed, Recipient: users[1]},
	}
	repo := &mockRepos{users: users, deliveries: deliveries}
	srv := services.NewManageDelivery(repo, repo)
	testCases := []struct {
		userId, deliveryId int
		success            bool
	}{
		{userId: 1, deliveryId: 1, success: true},
		{userId: 1, deliveryId: 2, success: false},
		{userId: 1, deliveryId: 3, success: true},
		{userId: 1, deliveryId: 4, success: false},
		{userId: 4, deliveryId: 4, success: false},
		{userId: 4, deliveryId: 1, success: false},
	}
	asrt := assert.New(t)
	for i, testCase := range testCases {
		t.Logf("case %d \n", i)
		d, e := srv.AssignToCourier(ctx, uint(testCase.userId), uint(testCase.deliveryId))
		if testCase.success {
			asrt.Nil(e)
			asrt.NotNil(d)
			asrt.Equal(d.Courier.Id, uint(testCase.userId))
		} else {
			asrt.NotNil(e)
			asrt.True(d == nil || d.Courier == nil || d.Courier.Id != uint(testCase.userId))
		}
	}
}

func TestManageDelivery_Complete(t *testing.T) {
	deliveries := map[uint]*entities.Delivery{
		1: {Id: 1, Status: valueobjects.Created},
		2: {Id: 2, Status: valueobjects.Canceled},
		3: {Id: 3, Status: valueobjects.Delivers},
		4: {Id: 4, Status: valueobjects.Completed},
	}
	repo := &mockRepos{deliveries: deliveries}
	srv := services.NewManageDelivery(repo, repo)
	testCases := []struct {
		deliveryId uint
		success    bool
	}{
		{deliveryId: 1, success: false},
		{deliveryId: 2, success: false},
		{deliveryId: 3, success: true},
		{deliveryId: 4, success: false},
	}
	asrt := assert.New(t)
	for i, testCase := range testCases {
		t.Logf("case %d", i)
		d, err := srv.Complete(ctx, testCase.deliveryId)
		if testCase.success {
			asrt.Nil(err)
			asrt.Equal(d.Status, valueobjects.Completed)
			asrt.LessOrEqual(time.Now().Sub(d.UpdatedAt).Milliseconds(), int64(1))
		} else {
			asrt.NotNil(err)
		}
	}

}
