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

func SearchIndex(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword != "" {
		c.Redirect(http.StatusFound, util.GetSearchUrl(keyword))
		return
	}
	c.Redirect(http.StatusFound, conf.Conf.Site.BaseUrl)
}

func Search(c *gin.Context) {
	keyword := c.Param("keyword")
	pageStr := c.Param("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if page == 1 {
		c.Redirect(http.StatusMovedPermanently, util.GetSearchUrl(keyword))
		return
	}
	if err != nil {
		page = 1
	}
	noteList, tagListMap, maxPage, err := service.NoteByKeyword(keyword, page, conf.Conf.Site.NotePerPage)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("search.html", gin.H{
		"hasPrevPage": page > 1,
		"prevPage":    page - 1,
		"hasNextPage": page < maxPage,
		"nextPage":    page + 1,
		"curPage":     page,
		"keyword":     keyword,
		"noteList":    noteList,
		"tagListMap":  tagListMap,
		"site":        conf.Conf.Site,
		"social":      conf.Conf.Social,
		"prod":        conf.Conf.Env.Prod,
		"haha":        util.HahaGenarate(),
	}))
}
