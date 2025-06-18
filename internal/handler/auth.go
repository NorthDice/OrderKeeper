package handler

import (
	"OrderKeeper/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required"`
}

type SignUpResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

const (
	ErrCodeValidation = "VALIDATION_ERROR"
	ErrCodeInternal   = "INTERNAL_ERROR"
)

func (h *Handler) signUp(c *gin.Context) {

	start := time.Now()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	h.logger.Info("signup request started",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
	)

	var input SignUpRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("validation failed",
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

	h.logger.Info("signup validation passed",
		zap.String("email", input.Email),
		zap.String("username", input.Username),
		zap.String("client_ip", clientIP),
	)

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	id, err := h.services.Authorization.CreateUser(c.Request.Context(), user)
	if err != nil {

		h.logger.Error("user creation failed",
			zap.String("email", input.Email),
			zap.String("username", input.Username),
			zap.String("client_ip", clientIP),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)

		statusCode := http.StatusInternalServerError
		errorCode := ErrCodeInternal
		errorMsg := "Failed to create user"

		c.JSON(statusCode, ErrorResponse{
			Error: errorMsg,
			Code:  errorCode,
		})
		return
	}

	h.logger.Info("user created successfully",
		zap.Int("user_id", id),
		zap.String("email", input.Email),
		zap.String("username", input.Username),
		zap.String("client_ip", clientIP),
		zap.Int("status_code", http.StatusOK),
		zap.Duration("total_duration", time.Since(start)),
	)

	c.JSON(http.StatusOK, SignUpResponse{
		ID:      id,
		Message: "User created successfully",
	})
}
func (h *Handler) signIn(c *gin.Context) {

}
