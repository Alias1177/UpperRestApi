package table

type User struct { // Определяем структуру для хранения данных пользователя
	Name     string `json:"name" validate:"required,min=3"`               // Имя пользователя (обязательно, минимум 3 символа)
	Email    string `json:"email" validate:"required,email"`              // Email пользователя (обязательно, должен соответствовать формату email)
	Password string `json:"password,omitempty" validate:"required,min=8"` // Пароль пользователя (обязательно, минимум 8 символов; будет опущен в JSON, если пустой)
	Age      int    `json:"age"`                                          // Возраст пользователя (необязательное поле)
}
