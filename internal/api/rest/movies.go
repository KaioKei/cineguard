package rest

import "github.com/gin-gonic/gin"

func List(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"movies": "OK"})
}
