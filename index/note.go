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
)

func Note(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	note, err := service.NoteById(id)
	if err != nil {
		glog.Error(err)
		ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	c.Render(http.StatusOK, render.NewRender("note.html", gin.H{
		"note": note,
		"site": conf.Conf.Site,
	}))
}
