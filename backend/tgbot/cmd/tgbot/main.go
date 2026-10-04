package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_http_server "github.com/wydentis/vykladna/shared/core/transport_http/server"
	core_tgbot_client "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/client"
	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
	core_tgbot_server "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/server"
	echo_transport_tgbot "github.com/wydentis/vykladna/tgbot/internal/feature/echo/transport/tgbot"
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

	tgClientCfg := core_tgbot_client.NewConfigMust()
	tgClient := core_tgbot_client.New(tgClientCfg)

	echoTransport := echo_transport_tgbot.NewEchoTransport(tgClient)

	router := core_tgbot_server.NewRouter()
	router.RegisterRoutes(
		echoTransport.Routes()...,
	)

	httpServerCfg := core_http_server.NewConfigMust()
	tgServerCfg := core_tgbot_server.NewConfigMust()
	tgServer := core_tgbot_server.NewServer(
		tgServerCfg,
		httpServerCfg,
		logger,
		core_tgbot_middleware.Trace(),
	)

	tgServer.RegisterRouters(router)

	if err := tgServer.Run(ctx); err != nil {
		logger.Fatal("telegram server run error: ", err)
	}
}
