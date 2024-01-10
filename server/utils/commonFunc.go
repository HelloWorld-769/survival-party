package utils

import (
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
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

func TableIsEmpty(tablename string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM " + tablename + ");"
	db.QueryExecutor(query, &exists)

	return exists
}

func GetUserData(userId string) (*model.User, error) {

	var user model.User
	query := "select * from users where user_id =?"
	err := db.QueryExecutor(query, &user, userId)
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func GetUserGameStatsData(userId string) (*model.UserGameStats, error) {

	var userGameStats model.UserGameStats
	query := "select * from user_game_stats where user_id =?"
	err := db.QueryExecutor(query, &userGameStats, userId)
	if err != nil {

		return nil, err
	}

	return &userGameStats, nil

}
