package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanbolat18/parcel/deliveries/pkg/http/request"
)

type ApiAuthProxyMiddleware struct {
	decMap request.DecoratorMap
}

func NewApiAuthProxyMiddleware(decMap request.DecoratorMap) *ApiAuthProxyMiddleware {
	return &ApiAuthProxyMiddleware{decMap: decMap}
}

func (a *ApiAuthProxyMiddleware) Proxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		dec := request.ApiAuthProxyDecorator(request.EmptyDecorator(), auth)
		_ = a.decMap.Store(ctx, dec)
		ctx.Next()
		_ = a.decMap.Store(ctx, dec)
	}
}
