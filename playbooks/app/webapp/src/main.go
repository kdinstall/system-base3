package main

import (
	"log"
	"os"
	"webapp/src/config"
	database "webapp/src/lib/database/sqlite"
)

func main() {
	// DB 接続確立 & 初期化
	initDatabase()

	// ルータ起動
	router := initRouter()
	addr := ":" + config.GetEnv().ServerPort
	log.Printf("Server starting on http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// initDatabase はスキーマ適用と初期データ投入を行う
func initDatabase() {
	db := database.GetDB()

	// schema.sql 適用
	schema, err := os.ReadFile("install/schema.sql")
	if err != nil {
		log.Fatalf("schema.sql の読み込みに失敗: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatalf("schema.sql の実行に失敗: %v", err)
	}
	log.Println("Schema applied.")

	// レコード数チェック（0 件のときのみシード投入）
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM users"); err != nil {
		log.Fatalf("COUNT クエリに失敗: %v", err)
	}

	if count == 0 {
		seed, err := os.ReadFile("install/seed.sql")
		if err != nil {
			log.Fatalf("seed.sql の読み込みに失敗: %v", err)
		}
		if _, err := db.Exec(string(seed)); err != nil {
			log.Fatalf("seed.sql の実行に失敗: %v", err)
		}
		log.Println("Seed data inserted.")
	} else {
		log.Printf("DB already has %d user(s). Skipping seed.", count)
	}
}
