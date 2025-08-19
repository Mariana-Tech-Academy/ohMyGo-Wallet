package handler

import (
	"net/http"
	"strconv"
	"vaqua/models"
	"vaqua/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service *services.TransactionService
}


func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var tx models.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uidAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	tx.UserID = uidAny.(uint)

	
	if tx.Type == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "transaction type is required"})
    return
}
if tx.Amount <= 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
    return
}


if tx.Status == "" {
    tx.Status = "pending"
}

if err := h.Service.CreateTransaction(&tx); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
    return
}
c.JSON(http.StatusCreated, tx)

}


func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	uidAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := uidAny.(uint)

	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transactions, err := h.Service.GetAllTransactions(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}
	if len(transactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no transactions found"})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
