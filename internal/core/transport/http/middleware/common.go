package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	core_logger "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/logger"
	core_http_response "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-toxagpark-fa47c8d8/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func CORS(allowedOriginsList []string) Middleware {
	allowedOrigins := make(map[string]struct{})

	for _, origin := range allowedOriginsList {
		allowedOrigins[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if _, ok := allowedOrigins[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)
			c.Header(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
		}

		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequestID() Middleware {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Header(requestIDHeader, requestID)
		c.Set(requestIDHeader, requestID)

		c.Next()
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(c *gin.Context) {
		requestID := c.GetString(requestIDHeader)

		l := log.With(
			zap.String("request_id", requestID),
			zap.String("url", c.Request.URL.String()),
		)

		ctx := core_logger.ToContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}

func Trace() Middleware {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := core_logger.FromContext(ctx)

		before := time.Now()

		log.Debug(
			">>> incoming HTTP request",
			zap.String("http_method", c.Request.Method),
			zap.Time("time", before.UTC()),
		)

		c.Next()

		since := time.Since(before)
		statusCode := c.Writer.Status()

		log.Debug(
			"<<< done HTTP request",
			zap.Int("status_code", statusCode),
			zap.Duration("latency", since),
		)
	}

}

func Panic() Middleware {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := core_logger.FromContext(ctx)

		responseHandler := core_http_response.NewHTTPResponseHandler(
			log,
			c,
		)

		defer func() {
			if p := recover(); p != nil {
				responseHandler.PanicResponse(
					p,
					"during handle HTTP request got unexpected panic",
				)
			}
		}()

		c.Next()
	}

}
