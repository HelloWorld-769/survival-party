package utils

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsEmail(e string) bool {
	//e = strings.ToLower(e)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
func IsPassValid(password string) error {

	if len(password) < 8 {
		return errors.New("password is too short")

	}
	hasUpperCase := false
	hasLowerCase := false
	hasNumbers := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		} else if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		} else if char >= '0' && char <= '9' {
			hasNumbers = true
		} else if char >= '!' && char <= '/' {
			hasSpecial = true
		} else if char >= ':' && char <= '@' {
			hasSpecial = true
		}
	}

	if !hasUpperCase {
		return errors.New("password do not contain upperCase Character")
	}

	if !hasLowerCase {
		return errors.New("password do not contain LowerCase Character")
	}

	if !hasNumbers {
		return errors.New("password do not contain any numbers")
	}

	if !hasSpecial {
		return errors.New("password do not contain any special character")
	}
	return nil
}

func HashPassword(password string) (*string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bs)
	return &hashedPassword, nil
}

func CheckPasswordHash(password, hash string) bool {

	fmt.Println("inside password check ")
	fmt.Println("password hash", password, hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
