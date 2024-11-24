package storage

import (
	"errors"
	"fmt"
	"github.com/EthanGuo-coder/llm-backend-api/config"
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

// InitializeSQLite 初始化 SQLite 数据库
func InitializeSQLite() error {
	dbPath := config.AppConfig.SQLite.Path

	// 确保父目录存在
	dir := filepath.Dir(dbPath)
	if err := ensureDirectoryExists(dir); err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	// 打开数据库
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	// 设置数据库连接属性
	db.SetMaxOpenConns(10)   // 最大连接数
	db.SetMaxIdleConns(5)    // 最大空闲连接数
	db.SetConnMaxLifetime(0) // 禁止自动关闭连接

	// 初始化表
	if err := createTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	fmt.Println("SQLite initialized successfully!")
	return nil
}

// ensureDirectoryExists 确保目录存在
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// createTables 创建所需的表
func createTables(db *sql.DB) error {
	tableSchemas := []string{
		`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );
        `,
		`
        CREATE TABLE IF NOT EXISTS conversations (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
        );
        `,
	}

	for _, schema := range tableSchemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	}
	return nil
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

// FetchConversationsByUserID 从数据库中获取指定用户的所有会话
func FetchConversationsByUserID(userID int64) ([]models.Conversation, error) {
	db := GetDB()
	query := `
        SELECT id, title 
        FROM conversations 
        WHERE user_id = ? 
        ORDER BY ROWID DESC;
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, errors.New("failed to fetch conversations: " + err.Error())
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var conversation models.Conversation
		if err := rows.Scan(&conversation.ID, &conversation.Title); err != nil {
			return nil, errors.New("failed to scan conversation: " + err.Error())
		}
		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("row iteration error: " + err.Error())
	}

	return conversations, nil
}
