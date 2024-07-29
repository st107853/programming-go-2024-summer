package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pet-progect.com/album"
)

func main() {
	album.Connect()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON
func getAlbums(c *gin.Context) {
	alb, err := album.Albums()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, alb)
		return
	}
	c.IndentedJSON(http.StatusOK, alb)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album.Album

	//Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	//Add the new album to the slice.
	id, err := album.AddAlbum(newAlbum)

	if err != nil {
		c.IndentedJSON(http.StatusNotExtended, "")
		return
	}
	mes := "new album id: " + strconv.Itoa(int(id))
	c.IndentedJSON(http.StatusCreated, gin.H{"message": mes})
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	//loop over the list of albums, looking for
	//an album whose ID value matchea the parameter.
	alb, err := album.AlbumByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, alb)
}
