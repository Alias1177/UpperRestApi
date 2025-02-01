package chain

import (
	"UpperRestApi/handler"                // Импортируем пакет с обработчиками
	"github.com/go-chi/chi/v5"            // Импортируем роутер chi
	"github.com/go-chi/chi/v5/middleware" // Импортируем встроенные middlewares
)

// AttachRouts подключает маршруты к роутеру через объект Handler
func AttachRouts(r chi.Router, h *handler.Handler) {
	r.Use(middleware.Recoverer) // Добавляем middleware для восстановления после паник
	r.Use(middleware.Logger)    // Добавляем middleware для логирования всех запросов

	r.Route("/api", func(r chi.Router) { // Группируем маршруты с префиксом /api
		r.Post("/", h.PostCreateUser)                 // Маршрут для создания нового пользователя
		r.Get("/{id}", h.GetUserByID)                 // Маршрут для получения пользователя по ID
		r.Delete("/{id}", h.DeleteUserByID)           // Маршрут для удаления пользователя по ID
		r.Put("/update/{login}", h.UpdateUserByLogin) // Маршрут для обновления данных пользователя по логину
	})
}
