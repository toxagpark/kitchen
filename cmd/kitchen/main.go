package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/config"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	core_pgx_pool "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/middleware"
	core_http_server "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/server"
	menu_pg_repository "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/menu/repository/postgres"
	menu_service "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/menu/service"
	menu_transport_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/menu/transport/http"
	order_pg_repository "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order/repository/postgres"
	order_service "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order/service"
	order_transport_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/order/transport/http"
	restaurant_pg_repository "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant/repository/postgres"
	restaurant_service "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant/service"
	restaurant_transport_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant/transport/http"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/features/restaurant_client"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Hello, KITCHEN!!!")

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		log.Fatalf("failed to init application logger: %s", err.Error())
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "restaurants"))
	restaurantRepo := restaurant_pg_repository.NewRepository(pool)
	restaurantsService := restaurant_service.NewService(restaurantRepo)
	restaurantTransport := restaurant_transport_http.NewRestaurantHTTPHandler(restaurantsService)

	logger.Debug("initializing feature", zap.String("feature", "menu"))
	menuRepo := menu_pg_repository.NewRepository(pool)
	menuService := menu_service.NewService(menuRepo)
	menuTransport := menu_transport_http.NewRestaurantHTTPHandler(menuService)

	logger.Debug("initializing feature", zap.String("feature", "order"))
	orderRepo := order_pg_repository.NewRepository(pool)
	restClientCfg := restaurant_client.NewConfigMust()
	restClient := restaurant_client.NewClient(restClientCfg, http.DefaultClient)
	orderService := order_service.NewService(pool, orderRepo, restClient)
	orderTransport := order_transport_http.NewOrderHTTPHandler(orderService)

	// logger.Debug("initializing feature", zap.String("feature", "tasks"))
	// tasksRepo := tasks_postgres_repository.NewTasksRepository(pool)
	// tasksService := tasks_service.NewTasksService(tasksRepo)
	// tasksTransport := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("starting kitchen application")

	httpCfg := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpCfg,
		logger,
		core_http_middleware.CORS(httpCfg.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	group := httpServer.Router.Group("/api/v1")

	restaurantTransport.Routes(group)
	menuTransport.Routes(group)
	orderTransport.Routes(group)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
