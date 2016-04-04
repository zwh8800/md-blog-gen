package index

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
)

func Index(c *gin.Context) {
	const limit = 10
	pageStr := c.Param("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	noteList, tagListMap, maxPage, err := service.NotesOrderByTime(page, limit)
	if err != nil {
		errorHandler(c, http.StatusInternalServerError, err)
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
	}))
}
