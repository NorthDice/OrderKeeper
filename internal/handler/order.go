package handler

import (
	"OrderKeeper/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type CreateOrderRequest struct {
	Status models.OrderStatus `json:"status"`
}

type CreateOrderResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
type GetOrdersResponse struct {
	Orders  []models.Order `json:"orders"`
	Message string         `json:"message"`
}

type DeleteOrderResponse struct {
	Message string `json:"message"`
}

type UpdateOrderResponse struct {
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

func (h *Handler) getOrders(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("get orders request started",
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

	orders, err := h.services.Order.GetOrders(c.Request.Context(), userId)
	if err != nil {
		h.logger.Error("failed to get orders",
			zap.Int("user_id", userId),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get orders",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("orders retrieved successfully",
		zap.Int("user_id", userId),
		zap.Int("orders_count", len(orders)),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusOK),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusOK, GetOrdersResponse{
		Orders:  orders,
		Message: "Orders retrieved successfully",
	})
}
func (h *Handler) getOrderById(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("get order by id request started",
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

	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid order id",
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Code:    ErrCodeValidation,
			Details: err.Error(),
		})
		return
	}

	order, err := h.services.Order.GetOrderByID(c.Request.Context(), userId, orderId)
	if err != nil {
		h.logger.Error("failed to get order by id",
			zap.Int("user_id", userId),
			zap.Int("order_id", orderId),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get order",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("order retrieved successfully",
		zap.Int("user_id", userId),
		zap.Int("order_id", order.ID),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusOK),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusOK, order)
}
func (h *Handler) updateOrder(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("update order request started",
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
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid order id",
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Code:    ErrCodeValidation,
			Details: err.Error(),
		})
		return
	}
	var input models.OrderUpdateInput
	if err = c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("failed to bind json",
			zap.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid input data",
			Code:    ErrCodeValidation,
			Details: err.Error(),
		})
	}
	if err = h.services.Order.UpdateOrder(c.Request.Context(), userId, orderId, input); err != nil {
		h.logger.Error("failed to update order",
			zap.Int("user_id", userId),
			zap.Int("order_id", orderId),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update order",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})
		return
	}
	h.logger.Info("order updated successfully",
		zap.Int("user_id", userId),
		zap.Int("order_id", orderId),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusOK),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusOK, UpdateOrderResponse{
		Message: "Order updated successfully",
	})
}
func (h *Handler) deleteOrder(c *gin.Context) {
	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("delete order request started",
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

	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("invalid order id",
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid order ID",
			Code:    ErrCodeValidation,
			Details: err.Error(),
		})
		return
	}

	if err = h.services.Order.DeleteOrder(c.Request.Context(), userId, orderId); err != nil {
		h.logger.Error("failed to delete order",
			zap.Int("user_id", userId),
			zap.Int("order_id", orderId),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete order",
			Code:    ErrCodeInternal,
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("order deleted successfully",
		zap.Int("user_id", userId),
		zap.Int("order_id", orderId),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusOK),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusOK, DeleteOrderResponse{
		Message: "Order deleted successfully",
	})
}
