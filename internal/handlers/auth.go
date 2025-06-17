package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"success": true,
	})
}
func (h *Handler) signIn(c *gin.Context) {

}
