package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// getEnv membaca environment variable dengan key, dan mengembalikan defaultVal jika belum di-set
func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func Init() {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_NAME", "mdl")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal koneksi ke DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Tidak bisa menghubungi DB: %v", err)
	}

	log.Println("Berhasil terhubung ke DB Moodle MySQL")
}

func Close() {
	if DB != nil {
		_ = DB.Close()
	}
}
