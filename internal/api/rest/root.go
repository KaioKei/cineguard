package rest

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"status": "OK"})
}
