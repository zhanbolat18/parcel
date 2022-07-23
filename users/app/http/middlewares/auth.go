package middlewares

import (
	"github.com/gin-gonic/gin"
	httpLib "github.com/zhanbolat18/parcel/libs/http"
	"github.com/zhanbolat18/parcel/users/internal/services"
	"net/http"
)

const BearerSchema = "Bearer "
const AuthHeader = "Authorization"

type AuthMiddleware struct {
	srv *services.AuthService
}

func NewAuthMiddleware(srv *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{srv: srv}
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(AuthHeader)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
			return
		}
		token := header[len(BearerSchema):]
		u, err := a.srv.Authorization(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpLib.Unauthorized())
			return
		}
		ctx.Set("user", u)
		ctx.Next()
	}
}
