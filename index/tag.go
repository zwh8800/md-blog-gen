package index

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
)

func AllTag(c *gin.Context) {
	tagList, noteListMap, err := service.AllNotesTags()
	if err != nil {
		errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Render(http.StatusOK, render.NewRender("allTag.html", gin.H{
		"tagList":     tagList,
		"noteListMap": noteListMap,
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
	if useId {
		tag, noteList, err = service.NotesByTagId(tagId)
	} else {
		tag, noteList, err = service.NotesByTagName(tagName)
	}
	if err != nil {
		errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	c.Render(http.StatusOK, render.NewRender("tag.html", gin.H{
		"tag":      tag,
		"noteList": noteList,
		"site":     conf.Conf.Site,
	}))
}
