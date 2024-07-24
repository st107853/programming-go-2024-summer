package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"

	//	"github.com/go-gomail/gomail"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	app.Get("/profile/:guid", Profile)
	app.Post("/refresh/:guid/:token", Refresh)

	logrus.Fatal(app.Listen(":8080"))

	//Capture connection propeties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "guids",
	}
	//Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

var db *sql.DB

type User struct {
	GUID  string
	Email string
	IP    string
}

// Handler of http requests for token updates
func Refresh(c *fiber.Ctx) error {
	var user User

	userGuid := c.Params("guid")

	if userGuid == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//Getting data from MySql
	row := db.QueryRow("SELECT * FROM users WHERE guids = ?", userGuid)

	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		return err
	}

	//Encoding the refresh token
	base64Token := c.Params("token")
	token, err := base64.StdEncoding.DecodeString(base64Token)
	if err != nil {
		return err
	}

	if string(token) != refreshToken(user) {
		return errors.New("wrong refresh token")
	}
	/*
		//Sending an email from to the mail if the user's IP has changed
			if user.IP != string(c.Context().LocalIP()) {
				fmt.Println("user ip pu-pu-pu")
				user.IP = string(c.Context().LocalIP())
				m := gomail.NewMessage()
				m.SetHeader("From", "alex@example.com")
				m.SetHeader("To", user.Email)
				m.SetHeader("Subject", "Hello!")
				m.SetBody("text/html", "Warning! You changed your ip and token.")

				d := gomail.NewDialer("test.com/m", 587, "user", "123456")

				//Send the email to User
				if err := d.DialAndSend(m); err != nil {
					panic(err)
				}
			}
	*/
	jwtSecretKey = randString(16)

	return nil
}

// The structure of the HTTP response with tokens information
type ProfileResponse struct {
	Access  string `json:"jwtAccess"`
	Refresh string `json:"jwtRefresh"`
}

// Handler for HTTP requests for user information
func Profile(c *fiber.Ctx) error {
	var user User

	userGuid := c.Params("guid")

	fmt.Println(userGuid)
	if userGuid == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//Getting data from MySql
	row := db.QueryRow("SELECT * FROM users WHERE guids = ?", userGuid)

	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		return err
	}

	base64Token := base64.StdEncoding.EncodeToString([]byte(refreshToken(user)))

	return c.JSON(ProfileResponse{
		Access:  jwtFromUser(user),
		Refresh: base64Token,
	})
}

// The secret key for signing the JWT token
var jwtSecretKey = []byte("very-secret-key")

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Updating the key after refresh
func randString(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return b
}

// Generating a JWT token
func jwtFromUser(user User) string {

	//Generating useful data that will be stored in the token
	payload := jwt.MapClaims{
		"sub": user.Email,
		"ip":  user.IP,
	}

	//Creating a new JWT token and signing it with the HS512 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return t
}

// Generating a JWT refresh token
func refreshToken(user User) string {

	//Generating useful data that will be stored in the token
	payload := jwt.MapClaims{
		"sub": user.GUID,
	}

	//Creating a new JWT token and signing it with the HS512 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return t
}
