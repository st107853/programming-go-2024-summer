package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v5"
)

var db *sql.DB

func main() {
	var err error
	//Capture connection propeties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "testtask",
	}
	//Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	var port = 8080

	mux := http.NewServeMux()

	mux.HandleFunc("GET /profile/{guid}", Profile)
	mux.HandleFunc("POST /refresh/{guid}/{token}", Refresh)

	// create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	log.Printf("Server started at :%v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	GUID  string
	Email string
	IP    string
}

// Handler of http requests for token updates
func Refresh(w http.ResponseWriter, r *http.Request) {
	var user User

	userGuid := r.PathValue("guid")
	if userGuid == "" {
		http.Error(w, "Empty guid", http.StatusBadRequest)
		return
	}

	//Getting data from MySql
	row := db.QueryRow("SELECT * FROM users WHERE guids = ?", userGuid)

	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//Encoding the refresh token
	base64Token := r.PathValue("token")
	token, err := base64.StdEncoding.DecodeString(base64Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if string(token) != refreshToken(user) {
		http.Error(w, "Wrong token", http.StatusBadRequest)
		return
	}

	fmt.Println(r.RemoteAddr)

	//Sending an email from to the mail if the user's IP has changed
	if user.IP != strings.Split(r.RemoteAddr, ":")[0] {
		fmt.Println("user ip pu-pu-pu")
		if err := emailSender(user.Email); err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println("some problem")
		}
	}

	jwtSecretKey = randString(16)

	w.Write([]byte("All done."))
}

// The structure of the HTTP response with tokens information
type ProfileResponse struct {
	Access  string `json:"jwtAccess"`
	Refresh string `json:"jwtRefresh"`
}

// Handler for HTTP requests for user information
func Profile(w http.ResponseWriter, r *http.Request) {
	var user User

	userGuid := r.PathValue("guid")
	if userGuid == "" {
		http.Error(w, "Empty guid", http.StatusBadRequest)
		return
	}

	//Getting data from MySql
	row := db.QueryRow("SELECT * FROM users WHERE guids = ?", userGuid)

	if err := row.Scan(&user.GUID, &user.Email, &user.IP); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	base64Token := base64.StdEncoding.EncodeToString([]byte(refreshToken(user)))

	out, err := json.Marshal(ProfileResponse{
		Access:  jwtFromUser(user),
		Refresh: base64Token,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out)
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

func emailSender(rcpt string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail("katelugovaua0@gmail.com"); err != nil {
		return err
	}

	if err := c.Rcpt(rcpt); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wc, "Warning! You changed your ip and token.")
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
