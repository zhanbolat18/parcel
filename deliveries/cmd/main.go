package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhanbolat18/parcel/deliveries/app/http/controllers"
	"github.com/zhanbolat18/parcel/deliveries/app/http/middlewares"
	"github.com/zhanbolat18/parcel/deliveries/config"
	_ "github.com/zhanbolat18/parcel/deliveries/docs"
	"github.com/zhanbolat18/parcel/deliveries/internal/repositories"
	httpRepository "github.com/zhanbolat18/parcel/deliveries/internal/repositories/http"
	"github.com/zhanbolat18/parcel/deliveries/internal/repositories/postgres"
	"github.com/zhanbolat18/parcel/deliveries/internal/services"
	"github.com/zhanbolat18/parcel/deliveries/pkg/http/request"
	"go.uber.org/dig"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// @title Parcel Delivery Service
// @version 1.0
// @description

// @contact.name Zhanbolat
// @contact.email zhanbolat.nurutdin@gmail.com

// @host localhost:8081
// @BasePath /
// @schemes http
func main() {
	c := dig.New()
	provideDependencies(c)
	mustWork(c.Invoke(func(engine *gin.Engine) {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}))
	mustWork(c.Invoke(func(engine *gin.Engine,
		controller *controllers.Delivery,
		roleMw *middlewares.RoleMiddleware,
		authProxyMw *middlewares.ApiAuthProxyMiddleware,
		authMw *middlewares.AuthMiddleware) {
		engine.POST("/deliveries", authMw.Auth(), roleMw.CheckRole("user"), controller.Create)
		engine.GET("/deliveries", authMw.Auth(), roleMw.CheckRole("admin", "courier"), controller.GetAllDeliveries)
		engine.GET("/deliveries/:id", authMw.Auth(), roleMw.CheckRole("admin", "courier"), controller.GetOneDelivery)
		engine.PUT("/deliveries/:id/complete", authMw.Auth(), roleMw.CheckRole("courier"), controller.CompleteDelivery)
		engine.POST("/deliveries/:id/courier/:courierId",
			authMw.Auth(),
			roleMw.CheckRole("admin"),
			authProxyMw.Proxy(),
			controller.AssignToCourier)
	}))

	mustWork(c.Invoke(func(server *http.Server) {
		go func() {
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf(err.Error())
			}
		}()
	}))
	gracefulShutdown(c)
}

func provideDependencies(container *dig.Container) {
	mustWork(container.Provide(config.NewConfig))
	mustWork(container.Provide(func(cfg *config.Config) *sqlx.DB {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			cfg.PgSQL.User, cfg.PgSQL.Password, cfg.PgSQL.DBName, cfg.PgSQL.Host, cfg.PgSQL.Port)
		return sqlx.MustConnect(cfg.PgSQL.Driver, connStr)
	}))

	mustWork(container.Provide(func(
		client *http.Client,
		cfg *config.Config,
		requestDecorator request.RequestDecorator,
	) repositories.UsersRepository {
		return httpRepository.NewUserRepository(client, cfg.Services.UsersBaseUrl, requestDecorator)
	}))
	mustWork(container.Provide(request.NewDecoratorMaps))
	mustWork(container.Provide(func(decoratorMap request.DecoratorMap) request.RequestDecorator {
		return decoratorMap
	}))
	mustWork(container.Provide(postgres.NewDeliveryRepository))
	mustWork(container.Provide(services.NewManageDelivery))
	mustWork(container.Provide(middlewares.NewRoleMiddleware))
	mustWork(container.Provide(middlewares.NewApiAuthProxyMiddleware))
	mustWork(container.Provide(func(client *http.Client, cfg *config.Config) *middlewares.AuthMiddleware {
		return middlewares.NewAuthMiddleware(client, cfg.Services.UsersBaseUrl)
	}))
	mustWork(container.Provide(controllers.NewDeliveryController))

	mustWork(container.Provide(func(cfg *config.Config) *http.Client {
		return &http.Client{
			Timeout: cfg.HttpClient.Timeout,
		}
	}))
	mustWork(container.Provide(func() *gin.Engine {
		engine := gin.Default()
		engine.Use(gin.Logger(), gin.Recovery())
		return engine
	}))
	mustWork(container.Provide(func(engine *gin.Engine, cfg *config.Config) *http.Server {
		return &http.Server{
			Addr:    cfg.Server.Port,
			Handler: engine,
		}
	}))

}

func mustWork(e error) {
	if e != nil {
		panic(e)
	}
}

func gracefulShutdown(c *dig.Container) {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
	mustWork(c.Invoke(func(server *http.Server, cfg *config.Config) {
		ctx, cf := context.WithTimeout(context.Background(), cfg.Server.ShutdownTime)
		defer cf()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("server shutting down with error %v", err)
		}
	}))
}
