package database

import (
	"database/sql"
	"log"
	"sync"
	"webapp/src/config"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// BaseDb は sqlx.DB を埋め込んだDB操作の基底構造体
type BaseDb struct {
	*sqlx.DB
}

// getDB は sync.OnceValue によりアプリ全体で DB 接続を1つだけ生成する
var getDB = sync.OnceValue(func() *sqlx.DB {
	dbPath := config.GetEnv().DatabasePath

	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		log.Fatalf("DB connection failed (path=%s): %v", dbPath, err)
	}

	// SQLite の WAL モードを有効化（並行読み取りの改善）
	db.MustExec("PRAGMA journal_mode=WAL;")
	db.MustExec("PRAGMA foreign_keys=ON;")

	log.Printf("DB connected: %s", dbPath)
	return db
})

// GetDB はグローバルな DB 接続を返す
func GetDB() *sqlx.DB {
	return getDB()
}

// NewBaseDb はグローバル DB 接続を使用した BaseDb を返す
func NewBaseDb() *BaseDb {
	return &BaseDb{DB: GetDB()}
}

// NewBaseDbWithDB はテスト用に任意の接続を注入した BaseDb を返す
func NewBaseDbWithDB(conn *sqlx.DB) *BaseDb {
	return &BaseDb{DB: conn}
}

// DoInTx はトランザクション内でコールバックを実行する
func (db *BaseDb) DoInTx(f func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// QueryRows はクエリを実行し、複数レコードを []map[string]interface{} で返す
func (db *BaseDb) QueryRows(query string, args ...interface{}) []map[string]interface{} {
	rs, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("QueryRows error: %v", err)
		return nil
	}
	defer rs.Close()

	col, err := rs.Columns()
	if err != nil {
		log.Printf("QueryRows Columns error: %v", err)
		return nil
	}

	typeVal, err := rs.ColumnTypes()
	if err != nil {
		log.Printf("QueryRows ColumnTypes error: %v", err)
		return nil
	}

	results := make([]map[string]interface{}, 0)

	for rs.Next() {
		colVar := make([]interface{}, len(col))
		for i := 0; i < len(col); i++ {
			setColVarType(&colVar, i, typeVal[i].DatabaseTypeName())
		}

		if err := rs.Scan(colVar...); err != nil {
			log.Printf("QueryRows Scan error: %v", err)
			return nil
		}

		result := make(map[string]interface{})
		for j := 0; j < len(col); j++ {
			setResultValue(&result, col[j], colVar[j], typeVal[j].DatabaseTypeName())
		}
		results = append(results, result)
	}

	if err := rs.Err(); err != nil {
		log.Printf("QueryRows rows error: %v", err)
		return nil
	}

	return results
}

// QueryRow はクエリを実行し、1レコードを map[string]interface{} で返す
func (db *BaseDb) QueryRow(query string, args ...interface{}) map[string]interface{} {
	rs, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("QueryRow error: %v", err)
		return nil
	}
	defer rs.Close()

	col, err := rs.Columns()
	if err != nil {
		log.Printf("QueryRow Columns error: %v", err)
		return nil
	}

	typeVal, err := rs.ColumnTypes()
	if err != nil {
		log.Printf("QueryRow ColumnTypes error: %v", err)
		return nil
	}

	var result map[string]interface{}

	if rs.Next() {
		colVar := make([]interface{}, len(col))
		for i := 0; i < len(col); i++ {
			setColVarType(&colVar, i, typeVal[i].DatabaseTypeName())
		}

		if err := rs.Scan(colVar...); err != nil {
			log.Printf("QueryRow Scan error: %v", err)
			return nil
		}

		result = make(map[string]interface{})
		for j := 0; j < len(col); j++ {
			setResultValue(&result, col[j], colVar[j], typeVal[j].DatabaseTypeName())
		}
	}

	if err := rs.Err(); err != nil {
		log.Printf("QueryRow rows error: %v", err)
		return nil
	}

	return result
}

func setColVarType(colVar *[]interface{}, i int, typeName string) {
	switch typeName {
	case "INTEGER":
		var s sql.NullInt64
		(*colVar)[i] = &s
	case "REAL":
		var s sql.NullFloat64
		(*colVar)[i] = &s
	case "TEXT":
		var s sql.NullString
		(*colVar)[i] = &s
	case "BLOB":
		var s sql.NullString
		(*colVar)[i] = &s
	default:
		var s interface{}
		(*colVar)[i] = &s
	}
}

func setResultValue(result *map[string]interface{}, index string, colVar interface{}, typeName string) {
	switch typeName {
	case "INTEGER":
		temp := *(colVar.(*sql.NullInt64))
		if temp.Valid {
			(*result)[index] = temp.Int64
		} else {
			(*result)[index] = nil
		}
	case "REAL":
		temp := *(colVar.(*sql.NullFloat64))
		if temp.Valid {
			(*result)[index] = temp.Float64
		} else {
			(*result)[index] = nil
		}
	case "TEXT":
		temp := *(colVar.(*sql.NullString))
		if temp.Valid {
			(*result)[index] = temp.String
		} else {
			(*result)[index] = nil
		}
	case "BLOB":
		temp := *(colVar.(*sql.NullString))
		if temp.Valid {
			(*result)[index] = temp.String
		} else {
			(*result)[index] = nil
		}
	default:
		(*result)[index] = colVar
	}
}
