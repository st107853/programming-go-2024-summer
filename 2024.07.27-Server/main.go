package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"pet-progect.com/album"
)

var port = os.Getenv("PORT")

func main() {
	album.Connect()

	balancConnecting()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /albums", getAlbums)
	mux.HandleFunc("GET /albums/{id}", getAlbumByID)
	mux.HandleFunc("POST /albums/{title}/{artist}/{price}", postAlbums)

	// create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
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

// getAlbums responds with the list of all albums as JSON
func getAlbums(w http.ResponseWriter, r *http.Request) {
	alb, err := album.Albums()
	if err != nil {
		slog.Error("Get album: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(alb)
	if err != nil {
		slog.Error("Get album: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	slog.Info("Get albums")
	w.Write(out)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(w http.ResponseWriter, r *http.Request) {
	price, _ := strconv.ParseFloat(r.PathValue("price"), 64)

	var newAlbum = album.Album{
		Title:  r.PathValue("title"),
		Artist: r.PathValue("artist"),
		Price:  price,
	}

	//Add the new album to the slice.
	id, err := album.AddAlbum(newAlbum)
	if err != nil {
		slog.Error("Post albun: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	slog.Info("Post albums")

	mes := "new album id: " + strconv.Itoa(int(id))
	w.Write([]byte(mes))
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	//loop over the list of albums, looking for
	//an album whose ID value matchea the parameter.
	alb, err := album.AlbumByID(id)
	if err != nil {
		slog.Error("Get album by id: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(alb)
	if err != nil {
		slog.Error("Get album by id: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	slog.Info("Get album by id")

	w.Write(out)
}

func balancConnecting() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "http://localhost:8080/", nil)
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
