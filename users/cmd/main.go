package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhanbolat18/parcel/users/app/http/controllers"
	"github.com/zhanbolat18/parcel/users/app/http/middlewares"
	"github.com/zhanbolat18/parcel/users/config"
	_ "github.com/zhanbolat18/parcel/users/docs"
	"github.com/zhanbolat18/parcel/users/internal/repositories/postgres"
	"github.com/zhanbolat18/parcel/users/internal/services"
	"github.com/zhanbolat18/parcel/users/internal/valueobjects"
	"github.com/zhanbolat18/parcel/users/pkg/crypto"
	"github.com/zhanbolat18/parcel/users/pkg/jwt"
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

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	c := dig.New()
	provideDependencies(c)
	mustWork(c.Invoke(func(engine *gin.Engine) {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}))
	mustWork(c.Invoke(func(engine *gin.Engine, c *controllers.AuthController, mw *middlewares.AuthMiddleware) {
		engine.POST("/auth", mw.Auth(), c.Auth)
		engine.POST("/login", c.Login)
	}))
	mustWork(c.Invoke(func(engine *gin.Engine, c *controllers.UserController,
		authMw *middlewares.AuthMiddleware, roleMw *middlewares.RoleMiddleware) {
		engine.POST("/signup", c.SignUp)
		engine.POST("/courier", authMw.Auth(), roleMw.CheckRole(valueobjects.Admin), c.CreateCourier)
		engine.GET("/couriers", authMw.Auth(), roleMw.CheckRole(valueobjects.Admin), c.Couriers)
		engine.GET("/couriers/:id", authMw.Auth(), roleMw.CheckRole(valueobjects.Admin), c.Courier)
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
	mustWork(container.Provide(func(cfg *config.Config) jwt.Jwt {
		return jwt.NewJwtManager(cfg.Jwt.Ttl, cfg.Jwt.BaseTimeDelta, cfg.Jwt.SignKey)
	}))
	mustWork(container.Provide(func(cfg *config.Config) crypto.PasswordHasher {
		return crypto.NewPasswordHasher(cfg.PasswordHasher.Cost)
	}))
	mustWork(container.Provide(func(cfg *config.Config) *sqlx.DB {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			cfg.PgSQL.User, cfg.PgSQL.Password, cfg.PgSQL.DBName, cfg.PgSQL.Host, cfg.PgSQL.Port)
		return sqlx.MustConnect(cfg.PgSQL.Driver, connStr)
	}))

	mustWork(container.Provide(postgres.NewUserRepository))
	mustWork(container.Provide(services.NewUserService))
	mustWork(container.Provide(services.NewAuthService))
	mustWork(container.Provide(controllers.NewAuthController))
	mustWork(container.Provide(controllers.NewUserController))
	mustWork(container.Provide(middlewares.NewAuthMiddleware))
	mustWork(container.Provide(middlewares.NewRoleMiddleware))

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
