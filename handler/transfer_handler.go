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
	// Bind JSON body to the TransferRequest model
	var request models.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set in middleware)
	userIDCon, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// type is asserted as uint 
	userID, ok := userIDCon.(uint) //check userID is of type uint
	if !ok {
		//try floating
		userIDFloat, ok := userIDCon.(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID has invalid type"})
			return
		}
		userID = uint(userIDFloat)
	}

	// call service
	err := h.Service.TransferFunds(userID, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}


