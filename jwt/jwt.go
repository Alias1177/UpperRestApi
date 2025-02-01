package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("your_secret_key") // Секретный ключ для подписи JWT

// GenerateHashedPassword генерирует хэшированный пароль с использованием JWT
func GenerateHashedPassword(password string) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute) // Устанавливаем время истечения действия токена

	claims := &jwt.RegisteredClaims{ // Создаем объект claims с данными
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Устанавливаем время истечения
		Subject:   password,                           // Используем пароль в качестве subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Создаем новый JWT токен с методом подписи HS256 и claims

	hashedPassword, err := token.SignedString(jwtKey) // Подписываем токен с использованием секретного ключа
	if err != nil {                                   // Если возникла ошибка при создании подписи
		return "", err // Возвращаем пустую строку и ошибку
	}

	return hashedPassword, nil // Возвращаем хэшированный пароль
}
