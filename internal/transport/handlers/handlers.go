package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Service interface {
	GetDataFromServiceLayerByID(ctx context.Context, id uint64) (string, error)
}

type Handler interface {
	Register(router *gin.Engine)
}

type handler struct {
	log     zerolog.Logger
	service Service
}

func (h *handler) Echo(c *gin.Context) {
	rawID := c.Query("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetDataFromServiceLayerByID(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, resp)
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/api/v1/echo", h.Echo)
}

func New(log zerolog.Logger, service Service) Handler {
	return &handler{
		log:     log,
		service: service,
	}
}
