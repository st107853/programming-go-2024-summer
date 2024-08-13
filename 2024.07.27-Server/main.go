package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"pet-progect.com/album"
)

func main() {
	album.Connect()

	var port = 8081

	mux := http.NewServeMux()

	mux.HandleFunc("GET /albums", getAlbums)
	mux.HandleFunc("GET /albums/{id}", getAlbumByID)
	mux.HandleFunc("POST /albums/{title}/{artist}/{price}", postAlbums)

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

// getAlbums responds with the list of all albums as JSON
func getAlbums(w http.ResponseWriter, r *http.Request) {
	alb, err := album.Albums()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(alb)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

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
		http.Error(w, err.Error(), 500)
		return
	}
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
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(alb)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out)
}
