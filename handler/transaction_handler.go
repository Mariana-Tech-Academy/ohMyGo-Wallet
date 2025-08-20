package handler

import (
	"net/http"
	"time"
	"vaqua/services"

	"github.com/gin-gonic/gin"
	"strconv"
) 
	


type TransactionHandler struct {
	Service *services.TransactionService
}

func (h *TransactionHandler) GetUserIncome(c *gin.Context) {
    userID := c.GetUint("userID") // from auth middleware

    // Placeholder times (no date filtering)
    var start, end time.Time

    income, _, err := h.Service.GetIncomeByPeriod(userID, start, end)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch income"})
        return
    }

    if len(income) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "no transactions found for this period"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"income": income})

}

func (h *TransactionHandler) GetUserExpenses( c*gin.Context) {
	userID := c.GetUint("userID") 

	var start, end time.Time

    expenses, _, err := h.Service.GetExpensesByPeriod(userID, start, end)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch expenses"})
        return
    }

    if len(expenses) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "no expenses found for this period"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"expenses": expenses})

}

func (h *TransactionHandler) GetBalance(c *gin.Context) {

    userID := c.GetUint("userID") 

    balance, err := h.Service.GetUserBalance(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not calculate balance"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"balance": balance})
}


//removed create transaction - this endpoint is not required given the transfer endpoint

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


func (h *TransactionHandler) GetTransaction(c *gin.Context) {

    user, ok := c.Get("user_id")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }
    userID := user.(uint)


    transaction, err := h.Service.GetTransactionByUserID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transaction"})
        return
    }


    if transaction.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "no transaction found"})
        return
    }

    c.JSON(http.StatusOK, transaction)
}

