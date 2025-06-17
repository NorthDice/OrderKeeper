package handlers

import (
	"OrderKeeper/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	services service.Service
	logger   *zap.Logger
}

func NewHandler(services service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	order := r.Group("/order", h.userIdentity)
	{
		order.POST("/", h.createOrder)
		order.GET("/orders", h.getOrders)
		order.GET("/", h.getOrderById)
		order.PUT("/:id", h.updateOrder)
		order.DELETE("/:id", h.deleteOrder)
	}

	return r
}
