package db

import (
	database "webapp/src/lib/database/sqlite"

	"github.com/jmoiron/sqlx"
)

// UserDb はユーザテーブルの CRUD を提供する
type UserDb struct {
	*database.BaseDb
}

// NewUserDb は UserDb を返す
func NewUserDb() *UserDb {
	return &UserDb{BaseDb: database.NewBaseDb()}
}

// ListUsers は全ユーザを取得する
func (d *UserDb) ListUsers() []map[string]interface{} {
	return d.QueryRows("SELECT id, name, email, created_at FROM users ORDER BY id ASC")
}

// GetUser は指定 ID のユーザを取得する
func (d *UserDb) GetUser(id int64) map[string]interface{} {
	return d.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id)
}

// CreateUser は新規ユーザを登録する
func (d *UserDb) CreateUser(name, email string) error {
	return d.DoInTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO users (name, email) VALUES (?, ?)",
			name, email,
		)
		return err
	})
}

// UpdateUser は指定 ID のユーザ情報を更新する
func (d *UserDb) UpdateUser(id int64, name, email string) error {
	return d.DoInTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec(
			"UPDATE users SET name = ?, email = ? WHERE id = ?",
			name, email, id,
		)
		return err
	})
}

// DeleteUser は指定 ID のユーザを削除する
func (d *UserDb) DeleteUser(id int64) error {
	return d.DoInTx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("DELETE FROM users WHERE id = ?", id)
		return err
	})
}
