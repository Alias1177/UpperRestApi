package connect

import (
	"github.com/jmoiron/sqlx" // Подключаем библиотеку для работы с базой данных
	_ "github.com/lib/pq"     // Импортируем драйвер для PostgreSQL
	"log"                     // Подключаем пакет для логирования
)

func Connect() *sqlx.DB {
	// Строка подключения к базе данных PostgreSQL
	dsn := "host=localhost port=5433 user=user password=password dbname=mydb sslmode=disable"

	// Устанавливаем соединение с базой данных с помощью sqlx
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil { // Если произошла ошибка при подключении
		// Логируем критическую ошибку и завершаем выполнение программы
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Возвращаем объект подключения к базе данных
	return db
}
