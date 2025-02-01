package repo

import (
	"UpperRestApi/handler"     // Импортируем пакет с обработчиком запросов
	"github.com/go-chi/chi/v5" // Импортируем библиотеку для работы с роутером
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jmoiron/sqlx" // Импортируем библиотеку для работы с базой данных
	"time"
)

type Repo struct { // Структура Repo инкапсулирует зависимости
	R      *chi.Mux         // Поле для роутера chi
	Db     *sqlx.DB         // Поле для подключения к базе данных
	Router *handler.Handler // Поле для обработчика запросов
}

// NewRepo создает новый объект Repo, который инкапсулирует все зависимости
func NewRepo(db *sqlx.DB) *Repo {
	r := chi.NewRouter()         // Создаем новый роутер chi
	h := handler.Handler{Db: db} // Инициализируем обработчик запросов, передавая подключение к базе данных
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.Limit(100, time.Minute)) // Ограничение до 100 запросов в минуту
	return &Repo{
		R:      r,  // Добавляем роутер chi в объект Repo
		Db:     db, // Добавляем подключение к базе данных в объект Repo
		Router: &h, // Устанавливаем ссылку на обработчик запросов в объект Repo
	}
}
