package rss

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/feeds"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func generateFeed() *feeds.Feed {
	noteList, _, _, err := service.NotesOrderByTime(1, conf.Conf.Site.NotePerPage)
	if err != nil {
		glog.Error(err)
		return nil
	}
	latestNoteUpdated := noteList[0].LastModified

	feed := &feeds.Feed{
		Title:       conf.Conf.Site.Name,
		Link:        &feeds.Link{Href: conf.Conf.Site.BaseUrl},
		Description: conf.Conf.Site.Name,
		Author:      &feeds.Author{Name: conf.Conf.Site.AuthorName, Email: conf.Conf.Site.AuthorEmail},
		Created:     latestNoteUpdated,
	}
	feed.Items = make([]*feeds.Item, 0)

	for _, note := range noteList {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          strconv.FormatInt(note.Id, 10),
			Title:       note.Title,
			Link:        &feeds.Link{Href: util.GetNoteUrl(note.Id)},
			Description: note.Content,
			Author:      &feeds.Author{Name: conf.Conf.Site.AuthorName, Email: conf.Conf.Site.AuthorEmail},
			Created:     note.Timestamp,
			Updated:     note.LastModified,
		})
	}
	return feed
}

func Rss(c *gin.Context) {
	feed := generateFeed()
	if feed == nil {
		index.ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}

	util.WriteContentType(c.Writer, "application/rss+xml; charset=utf-8")
	if err := feed.WriteRss(c.Writer); err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
}

func Atom(c *gin.Context) {
	feed := generateFeed()
	if feed == nil {
		index.ErrorHandler(c, http.StatusNotFound, errors.New("Not Found"))
		return
	}

	util.WriteContentType(c.Writer, "application/xml; charset=utf-8")
	if err := feed.WriteAtom(c.Writer); err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
}
