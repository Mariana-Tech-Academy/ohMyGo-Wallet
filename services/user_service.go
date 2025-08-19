package services

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"vaqua/middleware"
	"vaqua/models"
	"vaqua/redis"
	"vaqua/repository"
	"vaqua/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) SignUpNewUserAcct(newUser *models.SignUpRequest) error {
	// Check if email already exists
	_, err := s.Repo.GetUserByEmail(newUser.Email)
	if err == nil {
		return errors.New("this email is already in use")
	}

	// Hash password
	hashpass, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return err
	}
	newUser.Password = hashpass

	// Prepare new user
	user := &models.User{
		Email:    newUser.Email,
		Password: hashpass,
	}

	// Generate unique account number (string format)
	accNumStr := utils.GenerateRandomAccNumAsString()
	exists, err := s.Repo.CheckAccNumExists(accNumStr)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("account number already exists, please try again")
	}

	user.AccountNumber = accNumStr
	fmt.Printf("Saving new user: Email = %s | AccountNumber = %s\n", user.Email, user.AccountNumber)

	// Save user
	if err := s.Repo.CreateNewUser(user); err != nil {
		fmt.Println("Error saving user:", err)
		return err
	}

	// Create linked account record
	account := models.Account{
		UserID:        user.ID,
		AccountNumber: user.AccountNumber,
		Balance:       0,
	}
	if err := repository.NewTransferRepo().CreateAccount(&account); err != nil {
		fmt.Println("Error creating account for user:", err)
		return err
	}

	fmt.Println("User saved successfully to DB!")
	fmt.Println("Service: Finished SignUpNewUserAcct successfully")
	return nil
}


func (s *UserService) LoginUser(request models.LoginRequest) (string, error) {
    user, err := s.Repo.GetUserByEmail(request.Email)
    if err !=nil {
        return "", err
    }
    err = utils.ComparePassword(user.Password, request.Password)
    if err != nil {
        return "", err
    }
    token, err := middleware.GenerateJWT(user.ID, user.Email)
    if err != nil {
        return "", err
    }
    return token, nil 
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(email)
}

func (s *UserService) UpdateUserProfile(userID uint, updateUser *models.UpdateProfileRequest) (*models.User, error) {
	// get existing user
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// update required fields 
	if *updateUser.Firstname != "" {
		user.Firstname= updateUser.Firstname
	}
	if *updateUser.Lastname != "" {
		user.Lastname = updateUser.Lastname
	}
	if *updateUser.Phonenumber != 0 {
		user.Phonenumber = updateUser.Phonenumber
	}

	// save changes
	if err := s.Repo.UpdateUserProfile(user); err != nil {
		return nil, err
	}

	return user, nil
}
//LOG OUT USER
func (s *UserService) LogoutUser(c *gin.Context) error {
	
	//get token string from Auth header

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("authorization header required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// Check if token is already blacklisted
	blacklisted, err := redis.Client.Get(redis.Ctx, tokenString).Result()
	if err == nil && blacklisted == "blacklisted" {
		log.Println("Token is already blacklisted")
		return errors.New("token already blacklisted")
	}
	//parse the token get expiration time

	token, err := middleware.VerifyJWT(tokenString)
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	// using claims here to get expiration time so token is not expired before blacklisting
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("could not parse claims")
	}

	expUnix, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid  expiration time in token claims")
	}

	//store the token in redis with expiration time

	expirationTime := time.Unix(int64(expUnix), 0)
	duration := time.Until(expirationTime)

	err = redis.Client.Set(redis.Ctx, tokenString, "blacklisted", duration).Err()
	if err != nil {
		log.Printf("Failed to blacklist token from Redis: %v", err)
	} else {
		log.Println("Token successfully blacklisted in Redis")
	}

	//clear the token from context
	c.Set("token", nil)
	return nil
}
