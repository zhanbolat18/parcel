package controllers

import (
	"github.com/gin-gonic/gin"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"github.com/zhanbolat18/parcel/users/app/dto"
	"github.com/zhanbolat18/parcel/users/internal/services"
	"net/http"
	"strconv"
)

type UserController struct {
	srv *services.ManageUser
}

func NewUserController(srv *services.ManageUser) *UserController {
	return &UserController{srv: srv}
}

// SignUp godoc
// @Summary      SignUp
// @Description  sign up on service with email and password
// @Accept 		 json
// @Produce      json
// @Param        message  body  dto.UserDto  true  "courier info"
// @Success      200  {object}  entities.User
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Router       /signup [post]
func (u *UserController) SignUp(ctx *gin.Context) {
	uDto := &dto.UserDto{}
	err := ctx.ShouldBindJSON(uDto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := u.srv.SignUp(ctx, uDto.Email, uDto.Password)
	ctx.JSON(http.StatusOK, user)
}

// CreateCourier godoc
// @Summary      Create courier account
// @Description  Create courier account on service with email and password. Only admin have permission.
// @Accept 		 json
// @Produce      json
// @Param        message  body  dto.UserDto  true  "courier info"
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Success      200  {object}  entities.User
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /courier [post]
func (u *UserController) CreateCourier(ctx *gin.Context) {
	uDto := &dto.UserDto{}
	err := ctx.ShouldBindJSON(uDto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err))
		return
	}
	user, err := u.srv.CreateCourier(ctx, uDto.Email, uDto.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// Couriers godoc
// @Summary      Fetch all couriers account
// @Description  Fetch all couriers account on service with email and password. Only admin have permission.
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Success      200  {array}  []entities.User
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /couriers [get]
func (u *UserController) Couriers(ctx *gin.Context) {
	users, err := u.srv.Couriers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// Courier godoc
// @Summary      Fetch one courier account
// @Description  Fetch one courier account on service with email and password. Only admin have permission.
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Param		 id		path	integer	true	"courier id"
// @Success      200  {object}  entities.User
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /couriers/{id} [get]
func (u *UserController) Courier(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil || id < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest("invalid id"))
		return
	}
	users, err := u.srv.Courier(ctx, uint(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, users)
}
