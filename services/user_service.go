package services

import (
	"database/sql"
	"errors"

	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

func RegisterUser(username, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	db := storage.GetDB()
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return errors.New("username already exists")
	}

	return nil
}

func AuthenticateUser(username, password string) (string, error) {
	db := storage.GetDB()

	var user models.User
	row := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
