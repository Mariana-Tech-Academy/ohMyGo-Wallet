package handler

import (
	"net/http"
	"strconv"

	"vaqua/middleware"
	"vaqua/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

// GET /transactions/:id
func (h *TransactionHandler) GetTransactionByIDHandler(c *gin.Context) {
	//Extract transaction ID from URL
	idParam := c.Param("id")
	transactionID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	//Extract user ID from JWT token
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//transaction from service
	transaction, err := h.Service.GetTransactionByID(uint(transactionID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	//Return transaction details
	c.JSON(http.StatusOK, transaction)
}

// GET /dashboard/expenses
func (h *TransactionHandler) GetExpensesHandler(c *gin.Context) {
	//extract user from JWT token
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	//period query data
	period := c.DefaultQuery("period", "last_month")
	//get expenses from service
	expenses, err := h.Service.GetExpensesByUserAndPeriod(userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}

	//Return expenses
	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}
