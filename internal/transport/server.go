package transport

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler interface {
	Register(router *gin.Engine)
}

type Server interface {
	Run()
	Shutdown(ctx context.Context)
}

type server struct {
	srv      *http.Server
	log      zerolog.Logger
	host     string
	handlers Handler
}

func (s *server) Run() {
	r := gin.Default()
	s.handlers.Register(r)
	s.srv.Handler = r

	if err := s.srv.ListenAndServe(); err != nil {
		s.log.Fatal().Err(err).Send()
	}
}

func (s *server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Fatal().Err(err).Msg("failed to shutdown srv")
	}
}

func New(log zerolog.Logger, host string, handlers Handler) Server {
	s := &http.Server{
		Addr:         host,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	return &server{
		srv:      s,
		log:      log,
		host:     host,
		handlers: handlers,
	}
}
