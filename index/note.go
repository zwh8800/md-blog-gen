package index

import (
	"errors"
	"html/template"
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

var qrcodeCache = make(map[int64]string)

func Note(c *gin.Context) {
	useId := true

	notename := c.Param("id")
	id, err := strconv.ParseInt(notename, 10, 64)
	if err != nil {
		useId = false
	}
	var note *model.Note
	if useId {
		note, err = service.NoteById(id)
	} else {
		note, err = service.NoteByNotename(notename)
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
	if useId && note.Notename.Valid {
		c.Redirect(http.StatusMovedPermanently, util.GetNoteUrlByNotename(note.Notename.String))
	}
	// used by qrcode below
	id = note.Id

	qrcodeDataUrl, ok := qrcodeCache[id]
	if !ok {
		qrcodeDataUrl, err = util.GenerateQrcodePngDataUrl(util.GetNoteUrl(id))
		if err != nil {
			glog.Error(err)
		}
		qrcodeCache[id] = qrcodeDataUrl
	}

	relatedNotes, err := service.RelatedNote(note.Id)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	tags, err := service.TagsByNoteId(id)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	c.Render(http.StatusOK, render.NewRender("note.html", gin.H{
		"note":          note,
		"tags":          tags,
		"relatedNotes":  relatedNotes,
		"site":          conf.Conf.Site,
		"qrcodeDataUrl": template.URL(qrcodeDataUrl),
		"prod":          conf.Conf.Env.Prod,
	}))
}
