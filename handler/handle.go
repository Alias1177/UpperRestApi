package handler

import (
	"UpperRestApi/jwt"
	"UpperRestApi/table"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Handler struct {
	Db *sqlx.DB // Поле для хранения ссылки на базу данных
}

// Создать и занести данные пользователя в таблицу
func (h *Handler) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	var user table.User                                // Объявляем переменную для хранения данных пользователя
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок ответа в формате JSON

	// Декодируем тело запроса в структуру user
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest) // Отправляем ошибку, если формат JSON некорректен
		log.Printf("Ошибка при декодировании JSON: %v", err)      // Логируем ошибку
		return
	}

	validate := validator.New()  // Создаем объект валидатора
	err := validate.Struct(user) // Проверяем структуру user на соответствие правилам
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity) // Отправляем ошибку, если валидация не прошла
		return
	}

	if user.Password == "" { // Проверяем, указан ли пароль
		http.Error(w, "Пароль обязателен для регистрации", http.StatusBadRequest) // Отправляем ошибку, если пароль отсутствует
		log.Println("Ошибка: пароль обязателен")                                  // Логируем ошибку
		return
	}

	// Генерируем хэшированный пароль
	hashedPassword, err := jwt.GenerateHashedPassword(user.Password) // Хэшируем пароль
	if err != nil {                                                  // Если произошла ошибка при хэшировании
		http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError) // Отправляем ошибку
		log.Printf("Ошибка хэширования пароля: %v", err)                               // Логируем ошибку
		return
	}
	user.Password = hashedPassword // Сохраняем хэшированный пароль

	// Сохраняем пользователя в базу
	query := `INSERT INTO users (name, email, password, age) VALUES ($1, $2, $3, $4)` // SQL-запрос для вставки данных
	_, err = h.Db.Exec(query, user.Name, user.Email, user.Password, user.Age)         // Выполняем SQL-запрос
	if err != nil {                                                                   // Если произошла ошибка выполнения запроса
		http.Error(w, "Ошибка при сохранении данных в базу", http.StatusInternalServerError) // Отправляем ошибку
		log.Printf("Ошибка сохранения в БД: %v", err)                                        // Логируем ошибку
		return
	}

	w.WriteHeader(http.StatusCreated)                                                        // Устанавливаем статус ответа "Created"
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно добавлен"}) // Отправляем сообщение об успешном добавлении
}

// Получить пользователя по ID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")                        // Получаем ID пользователя из параметров URL
	var user table.User                                // Объявляем переменную для хранения данных пользователя
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок ответа в формате JSON

	// Запрос на получение данных пользователя
	query := `SELECT name, email, age FROM users WHERE id = $1` // SQL-запрос на выбор данных пользователя
	err := h.Db.Get(&user, query, id)                           // Выполняем запрос и сохраняем результат в переменной user
	if err != nil {                                             // Если пользователь не найден или возникла ошибка
		http.Error(w, "Пользователь не найден", http.StatusNotFound) // Отправляем ошибку
		log.Printf("Ошибка получения пользователя из БД: %v", err)   // Логируем ошибку
		return
	}

	w.WriteHeader(http.StatusOK)    // Устанавливаем статус ответа как "OK"
	json.NewEncoder(w).Encode(user) // Возвращаем данные пользователя в JSON-формате
}

// Удалить пользователя по ID
func (h *Handler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")                        // Получаем ID пользователя из параметров URL
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок ответа в формате JSON

	// Запрос на удаление пользователя
	query := `DELETE FROM users WHERE id = $1` // SQL-запрос на удаление пользователя
	result, err := h.Db.Exec(query, id)        // Выполняем запрос
	if err != nil {                            // Если запрос завершился с ошибкой
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError) // Отправляем ошибку
		log.Printf("Ошибка удаления из БД: %v", err)                                      // Логируем ошибку
		return
	}

	rowsAffected, err := result.RowsAffected() // Получаем количество затронутых строк
	if err != nil {                            // Если возникла ошибка при получении результата
		http.Error(w, "Ошибка получения результата удаления", http.StatusInternalServerError) // Отправляем ошибку
		log.Printf("Ошибка получения количества удаленных строк: %v", err)                    // Логируем ошибку
		return
	}

	if rowsAffected == 0 { // Если ни одна строка не была затронута
		http.Error(w, "Пользователь с таким ID не найден", http.StatusNotFound) // Отправляем ошибку
		return
	}

	w.WriteHeader(http.StatusOK)                                                           // Устанавливаем статус ответа как "OK"
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно удален"}) // Возвращаем сообщение об успешном удалении
}

// Обновить информацию пользователя по имени (логину)
func (h *Handler) UpdateUserByLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок ответа в формате JSON
	login := chi.URLParam(r, "login")                  // Получаем логин из параметров URL

	var user table.User                                           // Объявляем переменную для хранения новых данных пользователя
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil { // Декодируем тело запроса
		http.Error(w, "Некорректный JSON", http.StatusBadRequest) // Отправляем ошибку, если JSON некорректен
		log.Printf("Ошибка при декодировании JSON: %v", err)      // Логируем ошибку
		return
	}

	// Проверяем, существует ли пользователь с указанным логином
	var existingUser table.User                                        // Объявляем переменную для проверки существующего пользователя
	checkQuery := `SELECT name, email, age FROM users WHERE name = $1` // SQL-запрос на проверку существования пользователя
	err := h.Db.Get(&existingUser, checkQuery, login)                  // Выполняем запрос
	if err != nil {                                                    // Если пользователь не найден или запрос завершился ошибкой
		http.Error(w, "Пользователь не найден", http.StatusNotFound)        // Отправляем ошибку
		log.Printf("Пользователь с логином '%v' не найден: %v", login, err) // Логируем ошибку
		return
	}

	// Обновляем информацию о пользователе
	updateQuery := `UPDATE users SET email = $1, age = $2 WHERE name = $3` // SQL-запрос на обновление данных
	_, err = h.Db.Exec(updateQuery, user.Email, user.Age, login)           // Выполняем запрос с новыми данными
	if err != nil {                                                        // Если запрос завершился с ошибкой
		http.Error(w, "Ошибка при обновлении данных в базе", http.StatusInternalServerError) // Отправляем ошибку
		log.Printf("Ошибка обновления в БД: %v", err)                                        // Логируем ошибку
		return
	}

	w.WriteHeader(http.StatusOK)                                                                           // Устанавливаем статус ответа как "OK"
	json.NewEncoder(w).Encode(map[string]string{"message": "Информация о пользователе успешно обновлена"}) // Возвращаем сообщение об успешном обновлении
}
