package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	client       *http.Client
	usersAuthUrl string
}

func NewAuthMiddleware(client *http.Client, usersAuthUrl string) *AuthMiddleware {
	if client == nil {
		panic("client must be set")
	}
	return &AuthMiddleware{
		client:       client,
		usersAuthUrl: fmt.Sprintf("%s/auth", strings.TrimRight(usersAuthUrl, "/")),
	}
}

type userModel struct {
	Data struct {
		Id    uint   `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"data"`
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
			return
		}
		um, err := a.makeRequest(ctx, authHeader)
		fmt.Println(err)
		fmt.Println(a.usersAuthUrl)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized(err))
			return
		}
		ctx.Set("user", &entities.User{
			Id:    um.Data.Id,
			Email: um.Data.Email,
			Role:  um.Data.Role,
		})
		ctx.Next()
	}
}

func (a *AuthMiddleware) makeRequest(ctx context.Context, header string) (*userModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.usersAuthUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", header)
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("cannot check auth: %s", string(bytes)))
	}
	um := &userModel{}
	err = jsoniter.NewDecoder(res.Body).Decode(um)
	if err != nil {
		return nil, err
	}
	return um, nil
}
