package sitemap

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/controller/index"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"

	"github.com/chonglou/sitemap"
	"github.com/gin-gonic/gin"
)

func SiteMap(c *gin.Context) {
	s := sitemap.New()

	noteList, _, err := service.NotesWithoutTagOrderByTime(0, math.MaxInt64)
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	latestNoteUpdated := noteList[0].LastModified

	s.Add(&sitemap.Item{
		Link:    conf.Conf.Site.BaseUrl,
		Updated: latestNoteUpdated,
	})

	s.Add(&sitemap.Item{
		Link:    util.GetArchiveUrl(),
		Updated: latestNoteUpdated,
	})

	monthList, _, err := service.NoteGroupByMonth()
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}

	for _, month := range monthList {
		updated := time.Date(int(month.Year), time.Month(month.Month+1), -1, 0, 0, 0, 0, time.UTC)
		if latestNoteUpdated.Before(updated) {
			updated = latestNoteUpdated
		}

		s.Add(&sitemap.Item{
			Link:    util.GetArchiveMonthUrl(month.Year, month.Month),
			Updated: updated,
		})
	}

	for _, note := range noteList {
		if note.Notename.Valid {
			s.Add(&sitemap.Item{
				Link:    util.GetNoteUrlByNotename(note.Notename.String),
				Updated: note.LastModified,
			})
		}
		s.Add(&sitemap.Item{
			Link:    util.GetNoteUrl(note.Id),
			Updated: note.LastModified,
		})
	}

	tagList, err := service.TagsByCount()
	if err != nil {
		glog.Error(err)
		index.ErrorHandler(c, http.StatusServiceUnavailable, errors.New("Service Unavailable"))
		return
	}
	for _, tag := range tagList {
		s.Add(&sitemap.Item{
			Link:    util.GetTagUrl(tag.Id),
			Updated: latestNoteUpdated,
		})
		s.Add(&sitemap.Item{
			Link:    util.GetTagNameUrl(tag.Name),
			Updated: latestNoteUpdated,
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
			Updated: latestNoteUpdated,
		})
	}

	c.XML(http.StatusOK, s)
}
