package main

import (
	httpx "github.com/ViitoJooj/door/internal/http"
	"github.com/ViitoJooj/door/internal/http/handler"
	"github.com/ViitoJooj/door/internal/http/middlewares"
	"github.com/ViitoJooj/door/internal/repository"
	"github.com/ViitoJooj/door/internal/services"
	"github.com/ViitoJooj/door/pkg/database"
	"github.com/ViitoJooj/door/pkg/dotenv"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	dotenv.GetEnv()
	database.Conn()
	router := router.New()

	authRepo, applicationRepo := repository.NewSQLiteRepository(database.DB)
	authService := services.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	applicationService := services.NewApplicationService(applicationRepo, authRepo)
	applicationHandler := handler.NewApplicationHandler(applicationService)

	proxyService := services.NewProxyService()
	proxyHandler := handler.NewProxyHandler(proxyService)

	httpx.RegisterAuthRoutes(router, authHandler)
	httpx.RegisterApplicationRouters(router, applicationHandler)
	httpx.RegisterProxyRoutes(router, proxyHandler)

	handlerWithCors := middlewares.CorsMiddleware(router.Handler)

	fasthttp.ListenAndServe(":7171", handlerWithCors)
}
