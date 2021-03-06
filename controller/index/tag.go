package index

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func AllTag(c *gin.Context) {
	tagList, err := service.TagsByCount()
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("all_tag.html", gin.H{
		"tagList": tagList,
		"site":    conf.Conf.Site,
		"social":  conf.Conf.Social,
		"prod":    conf.Conf.Env.Prod,
		"haha":    util.HahaGenarate(),
	}))
}

func Tag(c *gin.Context) {
	useId := true

	tagName := c.Param("tag")
	tagId, err := strconv.ParseInt(tagName, 10, 64)
	if err != nil {
		useId = false
	}
	var tag *model.Tag
	var noteList []*model.Note
	var tagListMap map[int64][]*model.Tag
	if useId {
		tag, noteList, tagListMap, err = service.NotesByTagId(tagId)
	} else {
		tag, noteList, tagListMap, err = service.NotesByTagName(tagName)
	}
	if err != nil {
		glog.Error(err)
		if err == dbr.ErrNotFound {
			ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		} else {
			ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		}
		return
	}
	if useId {
		c.Redirect(http.StatusMovedPermanently, util.GetTagNameUrl(tag.Name))
		return
	}
	c.Render(http.StatusOK, render.NewRender("tag.html", gin.H{
		"tag":        tag,
		"noteList":   noteList,
		"tagListMap": tagListMap,
		"site":       conf.Conf.Site,
		"social":     conf.Conf.Social,
		"prod":       conf.Conf.Env.Prod,
		"haha":       util.HahaGenarate(),
	}))
}
