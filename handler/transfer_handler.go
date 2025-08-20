package handler

import (
	"net/http"
	"vaqua/models"
	"vaqua/services"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	Service services.TransferService
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	
	var request models.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	userIDCon, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	
	userID, ok := userIDCon.(uint) 
	if !ok {
		
		userIDFloat, ok := userIDCon.(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
			return
		}
		userID = uint(userIDFloat)
	}

	err := h.Service.TransferFunds(userID, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}


