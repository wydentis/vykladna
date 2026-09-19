package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	students_postgres_repository "github.com/wydentis/vykladna/backend/main/internal/feature/students/repository/postgres"
	students_service "github.com/wydentis/vykladna/backend/main/internal/feature/students/service"
	students_transport_http "github.com/wydentis/vykladna/backend/main/internal/feature/students/transport/http"
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

	studentsRepository := students_postgres_repository.NewStudentsRepository(pool)
	studentsService := students_service.NewStudentsService(studentsRepository)
	studentsTransportHTTP := students_transport_http.NewUsersHTTPTransport(studentsService)

	apiRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiRouterV1.RegisterRoutes(
		studentsTransportHTTP.Routes()...,
	)

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
		logger.Fatal("HTTP server run error: ", err)
	}
}
