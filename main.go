package main

import (
	"UpperRestApi/chain"
	"UpperRestApi/connect"
	"UpperRestApi/repo"
	"log"
	"net/http"
)

func main() {
	db := connect.Connect()                            // Подключение к базе данных
	defer db.Close()                                   // Закрытие базы данных при завершении
	repository := repo.NewRepo(db)                     // Создаем объект Repo
	chain.AttachRouts(repository.R, repository.Router) // Подключаем маршруты с помощью Repo

	log.Fatal(http.ListenAndServe(":8087", repository.R)) // Запуск сервера с роутером из Repo
}
