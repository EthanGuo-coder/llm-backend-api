package storage

// SQLiteSQL 用于存储 SQLite 的 SQL 语句
const (
	CreateTableUsers = `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );`

	CreateTableConversations = `
        CREATE TABLE IF NOT EXISTS conversations (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
        );`

	InsertConversation = `
        INSERT INTO conversations (id, title, user_id) 
        VALUES (?, ?, ?);`

	DeleteConversation = `
        DELETE FROM conversations 
        WHERE id = ? AND user_id = ?;`

	FetchConversations = `
        SELECT id, title 
        FROM conversations 
        WHERE user_id = ? 
        ORDER BY ROWID DESC;`
)
