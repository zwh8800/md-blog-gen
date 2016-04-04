package index

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zwh8800/md-blog-gen/render"
)

func errorHandler(c *gin.Context, status int, err error) {
	c.Render(http.StatusOK, render.NewRender("error.html", gin.H{
		"err": err,
	}))
}
