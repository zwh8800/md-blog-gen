package index

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
)

func AllTag(c *gin.Context) {
	tagList, noteListMap, tagListMap, err := service.AllNotesTags()
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusInternalServerError, errors.New("Service unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("all_tag.html", gin.H{
		"tagList":     tagList,
		"noteListMap": noteListMap,
		"tagListMap":  tagListMap,
		"site":        conf.Conf.Site,
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
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	c.Render(http.StatusOK, render.NewRender("tag.html", gin.H{
		"tag":        tag,
		"noteList":   noteList,
		"tagListMap": tagListMap,
		"site":       conf.Conf.Site,
	}))
}
