package main

import (
	"database/sql"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const (
	contextKeyUser = "ip"
)

func main() {
	app := fiber.New()

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: jwtSecretKey,
		},
		ContextKey: contextKeyUser,
	}))
	app.Get("/profile", Profile)
	app.Post("/refresh", Refresh)

	logrus.Fatal(app.Listen(":8080"))
}

var db *sql.DB

type User struct {
	GUID  string
	Email string
	IP    string
}

// Структура HTTP-ответа на вход в аккаунт
// В ответе содержится JWT-токен авторизованного пользователя
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// Секретный ключ для подписи JWT-токена
// Необходимо хранить в безопасном месте
var jwtSecretKey = []byte("very-secret-key")

// Обработчик HTTP-запросов на вход в аккаунт
func Refresh(c *fiber.Ctx) error {
	var user User
	userGuid := c.Context().Value("guid")
	//	refreshToken := c.Context().Value("token")

	// Ищем пользователя в памяти приложения по электронной почте
	row := db.QueryRow("SELECT * FROM users WHERE guid = ?", userGuid)
	// Если пользователь не найден, возвращаем ошибку
	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		return err
	}

	return nil
}

// Структура HTTP-ответа с информацией о пользователе
type ProfileResponse struct {
	JWT string `json:"jwt"`
}

func jwtFromUser(user User) string {
	// Генерируем JWT-токен для пользователя,
	// который он будет использовать в будущих HTTP-запросах

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": user.Email,
		"ip":  user.IP,
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return t
}

// Обработчик HTTP-запросов на получение информации о пользователе
func Profile(c *fiber.Ctx) error {
	var user User

	userGuid := c.Params("guid", "")
	if userGuid == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	row := db.QueryRow("SELECT * FROM users WHERE guid = ?", userGuid)

	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		return err
	}

	return c.JSON(ProfileResponse{
		JWT: jwtFromUser(user),
	})
}
