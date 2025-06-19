package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	IsEmptyString = ""
)
const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)
const (
	InvalidHeader = "INVALID_HEADER"
	InvalidToken  = "INVALID_TOKEN"
	EmptyToken    = "EMPTY_TOKEN"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == IsEmptyString {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "token is missing",
			Code:    EmptyToken,
			Details: "token is missing",
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid authorization header",
			Code:    InvalidToken,
			Details: "Invalid authorization header",
		})
		return
	}

	userId, err := h.services.Authorization.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid authorization header",
			Code:    InvalidHeader,
			Details: err.Error(),
		})
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, fmt.Errorf("user id is not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, fmt.Errorf("user id is of invalid type")
	}

	return idInt, nil
}
