package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"github.com/zhanbolat18/parcel/users/app/dto"
	"github.com/zhanbolat18/parcel/users/internal/services"
	"net/http"
)

type AuthController struct {
	srv *services.AuthService
}

func NewAuthController(srv *services.AuthService) *AuthController {
	return &AuthController{srv: srv}
}

// Auth godoc
// @Summary      Authorization
// @Description  authorization on service with JWT token and return user info
// @Produce      json
// @Security 	 ApiKeyAuth
// @Param 		 Authorization  header    string  true  "Authentication header"
// @Success      200  {object}  entities.User
// @Failure      401  {object}  object{error=string}
// @Failure      500  {string}  Internal Server Error
// @Router       /auth [post]
func (a *AuthController) Auth(ctx *gin.Context) {
	u, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, httpLib.InternalServErr())
		return
	}
	ctx.JSON(http.StatusOK, httpLib.Success(u))
	return
}

// Login godoc
// @Summary      Authentication
// @Description  authentication on service with email and password
// @Accept 		 json
// @Produce      json
// @Param        message  body  dto.UserDto  true  "login info"
// @Success      200  {object}  object{token=string}
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /login [post]
func (a *AuthController) Login(ctx *gin.Context) {
	credDto := &dto.UserDto{}
	err := ctx.BindJSON(credDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("invalid format")})
		return
	}
	t, err := a.srv.Authentication(ctx, credDto.Email, credDto.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"token": t})
}
