package route

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/util"
	"github.com/zwh8800/md-blog-gen/rss"
)

func Route(r *gin.Engine) {
	r.GET("/", index.Index)
	r.GET(path.Join(util.GetPageBase(), ":page"), index.Index)
	r.GET(util.GetTagBase(), index.AllTag)
	r.GET(path.Join(util.GetTagBase(), ":tag"), index.Tag)
	r.GET(path.Join(util.GetNoteBase(), ":id"), index.Note)

	r.GET("/rss", rss.Rss)

	r.Static("/static", "./static")
}
