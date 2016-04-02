package route

import (
	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/index"
)

func Route(r *gin.Engine) {
	r.LoadHTMLGlob("template/*")

	r.GET("/", index.Index)
	r.GET("/page/:page", index.Index)
	r.GET("/tag/:tag", index.Tag)
	r.GET("/note/:id", index.Note)
}
