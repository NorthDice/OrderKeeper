package handler

import (
	"OrderKeeper/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CreateOrderRequest struct {
	Status models.OrderStatus `json:"status"`
}

type CreateOrderResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func (h *Handler) createOrder(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("create order request started",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
	)

	userId, err := getUserId(c)
	if err != nil {
		h.logger.Error("failed to get user id",
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})

		return
	}

	var input CreateOrderRequest

	if err = c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("order validation failed",
			zap.String("client_ip", clientIP),
			zap.String("error", err.Error()),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input data",
			Code:    ErrCodeValidation,
			Details: err.Error(),
		})
		return
	}
	order := &models.Order{
		UserID:    userId,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	h.logger.Info("order validation passed",
		zap.Int("user_id", userId),
		zap.String("status", string(input.Status)),
		zap.String("client_ip", clientIP),
	)

	if err = h.services.Order.CreateOrder(c.Request.Context(), userId, order); err != nil {
		h.logger.Error("order creation failed",
			zap.Int("user_id", userId),
			zap.String("status", string(input.Status)),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create order",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})
		return
	}
	h.logger.Info("order created successfully",
		zap.Int("user_id", userId),
		zap.String("status", string(input.Status)),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusCreated),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusCreated, CreateOrderResponse{
		ID:      order.ID,
		Message: "Order created successfully",
	})
}

func (h *Handler) getOrders(c *gin.Context)    {}
func (h *Handler) getOrderById(c *gin.Context) {}
func (h *Handler) updateOrder(c *gin.Context)  {}
func (h *Handler) deleteOrder(c *gin.Context)  {}
