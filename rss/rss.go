package rss

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func Rss(c *gin.Context) {
	feed := feeds.Feed{
		Title:       conf.Conf.Site.Name,
		Link:        &feeds.Link{Href: conf.Conf.Site.BaseUrl},
		Description: conf.Conf.Site.Name,
		Author:      &feeds.Author{Name: conf.Conf.Site.AuthorName, Email: conf.Conf.Site.AuthorEmail},
		Created:     time.Now(),
	}
	feed.Items = make([]*feeds.Item, 0)

	noteList, _, _, err := service.NotesOrderByTime(1, conf.Conf.Site.NotePerPage)
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}
	for _, note := range noteList {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       note.Title,
			Link:        &feeds.Link{Href: util.GetNoteUrl(note.Id)},
			Description: note.Content,
			Author:      &feeds.Author{Name: conf.Conf.Site.AuthorName, Email: conf.Conf.Site.AuthorEmail},
			Created:     note.Timestamp,
		})
	}

	util.WriteContentType(c.Writer, []string{"application/rss+xml; charset=utf-8"})
	if err := feed.WriteRss(c.Writer); err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
}
