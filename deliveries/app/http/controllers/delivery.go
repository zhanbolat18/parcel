package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhanbolat18/parcel/deliveries/app/dto"
	"github.com/zhanbolat18/parcel/deliveries/internal/entities"
	"github.com/zhanbolat18/parcel/deliveries/internal/services"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"net/http"
	"strconv"
)

type Delivery struct {
	srv *services.ManageDelivery
}

func NewDeliveryController(srv *services.ManageDelivery) *Delivery {
	return &Delivery{srv: srv}
}

// Create godoc
// @Summary      create delivery order
// @Description  create delivery order with required destination. Only user have permission.
// @Accept 		 json
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Param        message  body  dto.Destination  true  "destination info"
// @Success      200  {object}  entities.Delivery
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /deliveries [post]
func (d *Delivery) Create(ctx *gin.Context) {
	u, ok := d.getUser(ctx)
	if !ok {
		return
	}

	dest := &dto.Destination{}
	if err := ctx.ShouldBindJSON(dest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}
	delivery, err := d.srv.Create(ctx, u, dest.D)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// AssignToCourier godoc
// @Summary      assign to courier
// @Description  assign to courier the delivery order. Only admin have permission.
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Param 		 id  			path	integer	true	"delivery id"
// @Param 		 courierId  	path	integer	true	"courier id"
// @Success      200  {object}  entities.Delivery
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /deliveries/{id}/courier/{courierId} [post]
func (d *Delivery) AssignToCourier(ctx *gin.Context) {
	id, ok := d.getUintParam(ctx, "id")
	if !ok {
		return
	}
	courierId, ok := d.getUintParam(ctx, "courierId")
	if !ok {
		return
	}
	delivery, err := d.srv.AssignToCourier(ctx, id, courierId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// GetAllDeliveries godoc
// @Summary      fetch all deliveries
// @Description  Fetch all deliveries. Only admin have permission to see all.
// @Description  If endpoint called with courier, only assigned deliveries returned
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Success      200  {array}  []entities.Delivery
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /deliveries [get]
func (d *Delivery) GetAllDeliveries(ctx *gin.Context) {
	user, ok := d.getUser(ctx)
	if !ok {
		return
	}

	var deliveries []*entities.Delivery
	var err error

	switch user.Role {
	case "admin":
		deliveries, err = d.srv.GetAll(ctx)
	case "courier":
		deliveries, err = d.srv.GetAllByCourier(ctx, user)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpLib.InternalServErr(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, deliveries)
}

// GetOneDelivery godoc
// @Summary      get one delivery
// @Description  Get one delivery. Only admin have permission to see all.
// @Description  If endpoint called with courier, only assigned deliveries returned
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Param 		 id  			path	integer	true	"delivery id"
// @Success      200  {array}  entities.Delivery
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /deliveries/{id} [get]
func (d *Delivery) GetOneDelivery(ctx *gin.Context) {
	user, ok := d.getUser(ctx)
	if !ok {
		return
	}
	id, ok := d.getUintParam(ctx, "id")
	if !ok {
		return
	}

	delivery, err := d.srv.GetOne(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpLib.InternalServErr(err.Error()))
		return
	}

	if user.Role == "courier" && (delivery.CourierId == nil || *delivery.CourierId != user.Id) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, httpLib.Forbidden())
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// CompleteDelivery godoc
// @Summary      complete delivery
// @Description  Complete delivery. Only assigned courier have permission.
// @Produce      json
// @Param 		 Authorization  header    string  true  "Authentication header. Usage 'Bearer {token}'"
// @Param 		 id  			path	integer	true	"delivery id"
// @Success      200  {object}  entities.Delivery
// @Failure      400  {object}  object{error=string}
// @Failure      401  {object}  object{error=string}
// @Failure      403  {object}  object{error=string}
// @Router       /deliveries/{id}/complete [put]
func (d *Delivery) CompleteDelivery(ctx *gin.Context) {
	u, ok := d.getUser(ctx)
	if !ok {
		return
	}
	id, ok := d.getUintParam(ctx, "id")
	if !ok {
		return
	}

	delivery, err := d.srv.Complete(ctx, uint(id), u)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpLib.InternalServErr(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

func (d *Delivery) getUintParam(ctx *gin.Context, param string) (uint, bool) {
	idStr := ctx.Param(param)
	if idStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(fmt.Sprintf("invalid %s", param)))
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, httpLib.BadRequest(fmt.Sprintf("invalid %s", param)))
		return 0, false
	}
	return uint(id), true
}

func (d *Delivery) getUser(ctx *gin.Context) (*entities.User, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
		return nil, false
	}
	u, ok := user.(*entities.User)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
		return nil, false
	}
	return u, true
}
