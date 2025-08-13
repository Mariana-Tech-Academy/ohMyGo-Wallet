package handler

import (
	"net/http"
	"net/mail"
	"strconv"
	"vaqua/models"
	"vaqua/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) SignUpNewUserAcct(c *gin.Context) {
	var newUser models.SignUpRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if _, err := mail.ParseAddress(newUser.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if len(newUser.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password should contain minimum 6 characters"})
		return
	}

	// call the service layer
	err := h.Service.SignUpNewUserAcct(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// response
	c.JSON(http.StatusCreated, gin.H {
		"message": "User created successfully",
		"email": newUser.Email,
	})

}

func (h *UserHandler) LoginUser(c *gin.Context) {
    var request models.LoginRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }	
    token, err := h.Service.LoginUser(request)
	
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}


func(h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID1, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userID1

	var updateUser models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updatedUser, err := h.Service.UpdateUserProfile(userID.(uint), &updateUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": updatedUser,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
    // Get user ID from JWT claims (set by AuthMiddleware)
    tokenUserID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
        return
    }

    // Read the ID from the request (URL param or query)
    idParam := c.Param("id")
    if idParam == "" {
        idParam = c.Query("id")
    }

    if idParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }

    id, err := strconv.Atoi(idParam)
    if err != nil || id < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Authorization check — block if token ID doesn't match request ID
    if uint(id) != tokenUserID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access this user's data"})
        return
    }

    // Fetch the user
    user, err := h.Service.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.Password = "" // hide password
    c.JSON(http.StatusOK, user)
}


func (h *UserHandler) GetUserByEmail(c *gin.Context) {
    // Get email from JWT claims
    tokenEmail, exists := c.Get("email")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
        return
    }

    // Fetch user by email from DB
    user, err := h.Service.GetUserByEmail(tokenEmail.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.Password = "" // hide password
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) LogoutUser(c *gin.Context) {
    err := h.Service.LogoutUser(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}