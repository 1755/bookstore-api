package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1755/bookstore-api/api/routers/health"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type ServerConfig struct {
	Address string `validate:"required,ip"`
	Port    int    `validate:"required,min=1,max=65535"`
}

type Server struct {
	engine *gin.Engine
	config *ServerConfig
	srv    *http.Server
}

var ServerModule = wire.NewSet(
	NewServer,
)

func NewServer(
	config *ServerConfig,
	health *health.RouterBuilder,
) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())

	root := engine.Group("/")

	health.Build(root)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Address, config.Port),
		Handler: engine,
	}

	return &Server{engine, config, srv}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
