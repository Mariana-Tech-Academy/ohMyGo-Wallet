package utils

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashpass), nil
}

func ComparePassword(hashPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func GenerateRandomAccNum() uint64 {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)

	min := uint64(1000000000)
	max := uint64(9999999999)

	return min + uint64(rand.Intn(int(max - min)))
}

//string needed for account number

func GenerateRandomAccNumAsString() string {
	num := GenerateRandomAccNum()
	return fmt.Sprintf("%010d", num) // always 10 digits
}
