package middlewares

import (
	"github.com/gin-gonic/gin"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"github.com/zhanbolat18/parcel/users/internal/entities"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
	"net/http"
)

type RoleMiddleware struct {
}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (r *RoleMiddleware) CheckRole(role ...valueobjects.Role) gin.HandlerFunc {
	if len(role) == 0 {
		panic("roles must be set")
	}
	rolesMap := make(map[valueobjects.Role]struct{})
	for _, v := range role {
		rolesMap[v] = struct{}{}
	}
	return func(ctx *gin.Context) {
		u, exists := ctx.Get("user")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
			return
		}
		user := u.(*entities.User)

		if _, ok := rolesMap[user.Role]; !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, httpLib.Forbidden())
			return
		}
		ctx.Next()
	}
}
