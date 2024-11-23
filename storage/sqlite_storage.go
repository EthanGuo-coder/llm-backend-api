package storage

import (
	"errors"
	"fmt"
	"os"

	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/EthanGuo-coder/llm-backend-api/models"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func InitializeSQLite(dbPath string) {
	// 确保父目录存在
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create directory for database: %v", err))
		}
	}

	// 打开数据库
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to SQLite: %v", err))
	}

	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`

	createConversationTable := `
    CREATE TABLE IF NOT EXISTS conversations (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(fmt.Sprintf("Failed to create users table: %v", err))
	}
	_, err = db.Exec(createConversationTable)
	if err != nil {
		panic(fmt.Sprintf("Failed to create conversations table: %v", err))
	}
	fmt.Println("SQLite initialized successfully!")
}

// SaveConversationToDB 将会话保存到数据库
func SaveConversationToDB(userID int64, conversation *models.Conversation) error {
	db := GetDB()
	// 插入会话记录到数据库
	query := `
        INSERT INTO conversations (id, title, user_id) 
        VALUES (?, ?, ?);
    `
	_, err := db.Exec(query, conversation.ID, conversation.Title, userID)
	if err != nil {
		return errors.New("failed to insert conversation: " + err.Error())
	}
	return nil
}

// DeleteConversationFromDB 从数据库中删除会话元信息
func DeleteConversationFromDB(userID int64, conversationID string) error {
	db := GetDB()

	query := `
        DELETE FROM conversations 
        WHERE id = ? AND user_id = ?;
    `
	_, err := db.Exec(query, conversationID, userID)
	if err != nil {
		return errors.New("failed to delete conversation from database: " + err.Error())
	}

	return nil
}
