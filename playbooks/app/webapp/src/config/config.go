package config

import "os"

// Env はアプリケーション設定を保持する構造体
type Env struct {
	AppName      string
	ServerPort   string
	DatabasePath string
}

// GetEnv は設定を返す（環境変数でオーバーライド可能）
func GetEnv() Env {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./user.sqlite3"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return Env{
		AppName:      "ユーザ管理",
		ServerPort:   port,
		DatabasePath: dbPath,
	}
}
