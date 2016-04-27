package index

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func Index(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	noteList, tagListMap, maxPage, err := service.NotesOrderByTime(page, conf.Conf.Site.NotePerPage)
	if err != nil || len(noteList) == 0 {
		glog.Error(err)
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("index.html", gin.H{
		"hasPrevPage": page > 1,
		"prevPage":    page - 1,
		"hasNextPage": page < maxPage,
		"nextPage":    page + 1,
		"curPage":     page,
		"noteList":    noteList,
		"tagListMap":  tagListMap,
		"site":        conf.Conf.Site,
		"social":      conf.Conf.Social,
		"prod":        conf.Conf.Env.Prod,
		"haha":        util.HahaGenarate(),
	}))
}
