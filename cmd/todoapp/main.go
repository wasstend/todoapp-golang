package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	core_logger_zap "github.com/wasstend/todoapp-golang/internal/core/logger/zap"
	core_pgx_pool "github.com/wasstend/todoapp-golang/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/wasstend/todoapp-golang/internal/core/transport/http/middleware"
	core_http_server "github.com/wasstend/todoapp-golang/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/wasstend/todoapp-golang/internal/features/statistics/repository/postgres"
	statistics_service "github.com/wasstend/todoapp-golang/internal/features/statistics/service"
	statistics_transport_http "github.com/wasstend/todoapp-golang/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/wasstend/todoapp-golang/internal/features/tasks/repository/postgres"
	tasks_service "github.com/wasstend/todoapp-golang/internal/features/tasks/service"
	tasks_transport_http "github.com/wasstend/todoapp-golang/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/wasstend/todoapp-golang/internal/features/users/repository/postgres"
	users_service "github.com/wasstend/todoapp-golang/internal/features/users/service"
	users_transport_http "github.com/wasstend/todoapp-golang/internal/features/users/transport/http"
	web_fs_repository "github.com/wasstend/todoapp-golang/internal/features/web/repository/file_system"
	web_service "github.com/wasstend/todoapp-golang/internal/features/web/service"
	web_transport_http "github.com/wasstend/todoapp-golang/internal/features/web/transport/http"

	_ "github.com/wasstend/todoapp-golang/docs"
)

// @title 		Golang Todo API
// @version 	1.0
// @description Todo Application RestAPI scheme
// @host 		127.0.0.1:5050
// @BasePath 	/api/v1
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger_zap.NewLogger(core_logger_zap.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("connection pool failed", core_logger.Error(err))
	}

	defer pool.Close()

	logger.Debug("initializing feature", core_logger.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUserRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", core_logger.String("feature", "tasks"))

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", core_logger.String("feature", "statistics"))

	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature", core_logger.String("feature", "web"))

	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing http server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP Server run error", core_logger.Error(err))
	}
}
