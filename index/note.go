package index

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/service"
)

func Note(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	note, err := service.NoteById(id)
	if err != nil {
		errorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	c.HTML(http.StatusOK, "note.html", gin.H{
		"note": note,
	})
}
