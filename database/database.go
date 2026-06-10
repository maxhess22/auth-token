package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	var dsn string

	// Revisamos si estamos en App Engine comprobando si existe la variable de la instancia
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	if instanceConnectionName != "" {
		// Conexión por Socket Unix para App Engine
		dsn = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true&charset=utf8mb4",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			instanceConnectionName,
			os.Getenv("DB_NAME"),
		)
		fmt.Println("Conectando a MySQL vía Socket Unix en App Engine...")
	} else {
		// Conexión TCP original (para desarrollo local)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		fmt.Println("Conectando a MySQL vía TCP (Local)...")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al abrir la conexión: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la BD: ", err)
	}

	return db
}
