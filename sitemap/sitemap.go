package sitemap

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/index"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"

	"github.com/chonglou/sitemap"
	"github.com/gin-gonic/gin"
)

func SiteMap(c *gin.Context) {
	s := sitemap.New()
	s.Add(&sitemap.Item{
		Link:    conf.Conf.Site.BaseUrl,
		Updated: time.Now(),
	})

	s.Add(&sitemap.Item{
		Link:    util.GetArchiveUrl(),
		Updated: time.Now(),
	})

	monthList, _, err := service.NoteGroupByMonth()
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	for _, month := range monthList {
		s.Add(&sitemap.Item{
			Link:    util.GetArchiveMonthUrl(month),
			Updated: time.Now()),
		})
	}

	noteList, _, err := service.NotesWithoutTagOrderByTime(0, math.MaxInt64)
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	for _, note := range noteList {
		if note.Notename.Valid {
			s.Add(&sitemap.Item{
				Link:    util.GetNoteUrlByNotename(note.Notename.String),
				Updated: note.Timestamp,
			})
		}
		s.Add(&sitemap.Item{
			Link:    util.GetNoteUrl(note.Id),
			Updated: note.Timestamp,
		})
	}

	tagList, err := service.Tags()
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	for _, tag := range tagList {
		s.Add(&sitemap.Item{
			Link:    util.GetTagUrl(tag.Id),
			Updated: time.Now(),
		})
		s.Add(&sitemap.Item{
			Link:    util.GetTagNameUrl(tag.Name),
			Updated: time.Now(),
		})
	}

	_, maxPage, err := service.NotesWithoutTagOrderByTime(0, conf.Conf.Site.NotePerPage)
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	for i := int64(1); i <= maxPage; i++ {
		s.Add(&sitemap.Item{
			Link:    util.GetPageUrl(i),
			Updated: time.Now(),
		})
	}

	c.XML(http.StatusOK, s)
}
