package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var port = os.Getenv("PORT")
var logger TransactionLogger

func main() {

	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/albums/{id}", albumGetHandler).Methods("GET")
	r.HandleFunc("/albums/{title}/{artist}/{price}", albumPostHandler).Methods("POST")

	// create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      r,
	}

	//	BalancConnecting()

	log.Printf("Server started at :%v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func albumPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title, artist, price := vars["title"], vars["artist"], vars["price"]

	id, err := Post(title, artist, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.WritePut(id, title, artist, price)

	res := fmt.Sprintf("new album id: %v", id)
	w.Write([]byte(res))
	// w.WriteHeader(http.StatusCreated)
}

func albumGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	value, err := Get(id)

	if errors.Is(err, ErrorNoSuchId) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, _ := json.Marshal(&value)

	w.Write(res)
}

func initializeTransactionLog() error {
	var err error

	logger, err = NewPostgresTransactionLogger(PostgresDbParams{
		host:     "localhost",
		dbName:   os.Getenv("DBNAME"),
		user:     os.Getenv("DBUSER"),
		password: os.Getenv("DBPASS"),
	})

	if err != nil {
		return fmt.Errorf("failet to create transaction logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	count, ok, e := 0, true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = Delete(e.Id)
				count++
			case EventPut: // Got a PUT event!
				_, err = Post(e.Title, e.Artist, e.Prise)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	logger.Run()

	go func() {
		for err := range logger.Err() {
			log.Print(err)
		}
	}()

	return err
}

func BalancConnecting() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}

	str := fmt.Sprintf("http://localhost:%v", port)

	req.Header.Add("TODO", "Add me")
	req.Header.Add("serv", str)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		return
	}
}
