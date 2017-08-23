package index

import (
	"github.com/gin-gonic/gin"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/render"
)

func ErrorHandler(c *gin.Context, status int, err error) {
	c.Render(status, render.NewRender("error.html", gin.H{
		"err":  err,
		"site": conf.Conf.Site,
		"prod": conf.Conf.Env.Prod,
	}))
}
