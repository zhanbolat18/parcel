package http

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	"github.com/zhanbolat18/parcel/deliveries/pkg/http/request"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type UserRepository struct {
	client           *http.Client
	baseUrl          string
	requestDecorator request.RequestDecorator
}

type userModel struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type errModel struct {
	Err interface{} `json:"error"`
}

func NewUserRepository(client *http.Client, baseUrl string, requestDecorator request.RequestDecorator) *UserRepository {
	if client == nil {
		panic("http client must be set")
	}
	_, err := url.Parse(baseUrl)
	if err != nil {
		panic("invalid url")
	}

	return &UserRepository{
		client:           client,
		baseUrl:          strings.TrimRight(baseUrl, "/"),
		requestDecorator: requestDecorator,
	}
}

func (u *UserRepository) GetCourier(ctx context.Context, id uint) (*entities.User, error) {
	path := fmt.Sprintf("%s/couriers/%d", u.baseUrl, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	u.requestDecorator.Decorate(req)
	fmt.Println(u.requestDecorator)
	res, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, u.responseNotOk(res)
	}

	um := &userModel{}
	err = jsoniter.NewDecoder(res.Body).Decode(um)
	if err != nil {
		return nil, err
	}
	return &entities.User{
		Id:    um.Id,
		Email: um.Email,
		Role:  um.Role,
	}, nil
}

func (u *UserRepository) GetRecipient(ctx context.Context, id uint) (*entities.User, error) {
	path := fmt.Sprintf("%s/auth", u.baseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	u.requestDecorator.Decorate(req)
	res, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, u.responseNotOk(res)
	}

	um := &userModel{}
	err = jsoniter.NewDecoder(res.Body).Decode(um)
	if err != nil {
		return nil, err
	}
	if um.Id != id {
		return nil, errors.New(fmt.Sprintf("access to user with id \"%d\" forbidden", id))
	}
	return &entities.User{
		Id:    um.Id,
		Email: um.Email,
		Role:  um.Role,
	}, nil
}

func (u *UserRepository) responseNotOk(res *http.Response) error {
	if strings.Contains(
		strings.Join(res.Header[http.CanonicalHeaderKey("Content-Type")], ""),
		"application/json") {
		e := &errModel{}
		err := jsoniter.NewDecoder(res.Body).Decode(e)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("request failed: %v", e.Err))
	}
	bodyTxt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("unsuccess request: %s", string(bodyTxt)))
}
