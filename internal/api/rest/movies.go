package rest

import "github.com/gin-gonic/gin"

func ListMovies(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"movies": "OK"})
}

func CreateMovie(c *gin.Context) {

}
