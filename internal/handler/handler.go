package handler

import (
	"OrderKeeper/internal/handler/metrics"
	"OrderKeeper/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Handler struct {
	services *service.Service
	logger   *zap.Logger
}

func NewHandler(services *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(metrics.MetricsMiddleware())

	r.GET("/health", h.healthCheck)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	order := r.Group("/order", h.userIdentity)
	{
		order.POST("/", h.createOrder)
		order.GET("/", h.getOrders)
		order.GET("/:id", h.getOrderById)
		order.PUT("/:id", h.updateOrder)
		order.DELETE("/:id", h.deleteOrder)
	}

	return r
}

func (h *Handler) healthCheck(c *gin.Context) {
	h.logger.Debug("Health check requested")
	c.JSON(200, gin.H{
		"status":    "ok",
		"service":   "myapp",
		"timestamp": gin.H{},
	})
}
