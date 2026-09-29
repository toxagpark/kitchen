package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	Router *gin.Engine
	cfg    Config
	log    *core_logger.Logger
}

func NewHTTPServer(
	cfg Config,
	log *core_logger.Logger,
	middleware ...gin.HandlerFunc,
) *HTTPServer {
	router := gin.New()

	router.Use(middleware...)

	return &HTTPServer{
		Router: router,
		cfg:    cfg,
		log:    log,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.Router,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.cfg.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}

	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("close HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
