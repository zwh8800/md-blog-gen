package index

import "github.com/gin-gonic/gin"

func errorHandler(c *gin.Context, status int, err error) {
	c.HTML(status, "error.html", gin.H{
		"err": err,
	})
}
