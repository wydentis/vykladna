package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	auth_service "github.com/wydentis/vykladna/backend/main/internal/feature/auth/service"
	auth_transport_http "github.com/wydentis/vykladna/backend/main/internal/feature/auth/transport/http"
	students_postgres_repository "github.com/wydentis/vykladna/backend/main/internal/feature/students/repository/postgres"
	students_service "github.com/wydentis/vykladna/backend/main/internal/feature/students/service"
	students_transport_http "github.com/wydentis/vykladna/backend/main/internal/feature/students/transport/http"
	users_postgres_repository "github.com/wydentis/vykladna/backend/main/internal/feature/users/repository/postgres"
	users_service "github.com/wydentis/vykladna/backend/main/internal/feature/users/service"
	users_transport_http "github.com/wydentis/vykladna/backend/main/internal/feature/users/transport/http"
	core_auth "github.com/wydentis/vykladna/shared/core/auth"
	core_pgx_pool "github.com/wydentis/vykladna/shared/core/db/postgres/pool/pgx"
	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to start logger: ", err)
		os.Exit(1)
	}

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres pool: ", err)
	}

	authManager := core_auth.NewManager(core_auth.NewConfigMust())

	studentsRepository := students_postgres_repository.NewStudentsRepository(pool)
	studentsService := students_service.NewStudentsService(studentsRepository)
	studentsTransportHTTP := students_transport_http.NewUsersHTTPTransport(studentsService, authManager)

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPTransport(usersService)

	authService := auth_service.NewAuthService(usersRepository, *authManager)
	authTransportHTTP := auth_transport_http.NewAuthTransportHTTP(authService)

	apiRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiRouterV1.RegisterRoutes(studentsTransportHTTP.Routes()...)
	apiRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	httpServer.RegisterAPIVersionRouters(apiRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Fatal("HTTP server run error: ", "error", err)
	}
}
