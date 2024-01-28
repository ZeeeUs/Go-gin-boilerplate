package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler interface {
	Register(router *gin.Engine)
}

type Server struct {
	log      zerolog.Logger
	host     string
	handlers Handler
}

func (s *Server) Run() {
	r := gin.Default()
	s.handlers.Register(r)

	server := &http.Server{
		Addr:         s.host,
		Handler:      r,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	if err := server.ListenAndServe(); err != nil {
		s.log.Fatal().Err(err).Msg("")
	}
}

func (s *Server) Stop() {
	panic("implement me")
}

func New(log zerolog.Logger, host string, handlers Handler) *Server {
	return &Server{
		log:      log,
		host:     host,
		handlers: handlers,
	}
}
