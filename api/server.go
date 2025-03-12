package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/1755/bookstore-api/api/middlewares"
	"github.com/1755/bookstore-api/api/routers/health"
	"github.com/1755/bookstore-api/api/routers/v1/authorsv1"
	"github.com/1755/bookstore-api/api/routers/v1/booksv1"
	docs "github.com/1755/bookstore-api/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Address string `validate:"required,ip"`
	Port    int    `validate:"required,min=1,max=65535"`
}

type Server struct {
	logger        *zap.Logger
	engine        *gin.Engine
	config        *ServerConfig
	srv           *http.Server
	monitoringSrv *http.Server
}

var ServerModule = wire.NewSet(
	NewServer,
)

func NewServer(
	config *ServerConfig,
	monitoring *prometheus.Registry,
	logger *zap.Logger,
	health *health.RouterBuilder,
	booksv1 *booksv1.RouterBuilder,
	authorsv1 *authorsv1.RouterBuilder,
) *Server {
	docs.SwaggerInfo.Title = "Bookstore API"
	docs.SwaggerInfo.BasePath = "/v1"

	engine := gin.New()
	engine.Use(
		gin.Recovery(),
		middlewares.Logger(logger),
		middlewares.MetricsMiddleware(),
		cors.Default(), // todo: do not use it in production, for demo is ok
	)

	root := engine.Group("/")

	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	health.Build(root)
	booksv1.Build(root)
	authorsv1.Build(root)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Address, config.Port),
		Handler: engine,
	}

	monitoringEngine := gin.New()
	monitoringEngine.GET("/metrics", gin.WrapH(promhttp.HandlerFor(monitoring, promhttp.HandlerOpts{})))

	monitoringSrv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Address, 9999),
		Handler: monitoringEngine,
	}

	return &Server{logger, engine, config, srv, monitoringSrv}
}

func (s *Server) Run() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.logger.Info("Starting server...", zap.String("address", s.srv.Addr))
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error: %v", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.logger.Info("Starting monitoring server...", zap.String("address", s.monitoringSrv.Addr))
		if err := s.monitoringSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Monitoring server error: %v", zap.Error(err))
		}
	}()

	wg.Wait()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error("Server shutdown error: %v", zap.Error(err))
	}

	if err := s.monitoringSrv.Shutdown(ctx); err != nil {
		s.logger.Error("Monitoring server shutdown error: %v", zap.Error(err))
	}

	return nil
}
