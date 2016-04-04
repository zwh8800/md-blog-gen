package route

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/util"
)

func Route(r *gin.Engine) {
	r.LoadHTMLGlob("template/*")

	r.GET("/", index.Index)
	r.GET(path.Join(util.GetPageBase(), ":page"), index.Index)
	r.GET(util.GetTagBase(), index.AllTag)
	r.GET(path.Join(util.GetTagBase(), ":tag"), index.Tag)
	r.GET(path.Join(util.GetNoteBase(), ":id"), index.Note)

	r.Static("/static", "./static")
}
