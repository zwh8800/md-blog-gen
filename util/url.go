package util

import (
	"net/url"
	"path"
	"strconv"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/conf"
)

func GetSiteDomain() string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.Host
}

func UrlJoin(subUrl string) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	u, err := url.Parse(subUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func UrlGetPath(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return u.Path
}

func GetNoteBase() string {
	return "/" + conf.Conf.Site.NoteUrl
}

func GetNoteUrl(id int64) string {
	return UrlJoin(path.Join(GetNoteBase(), strconv.FormatInt(id, 10)))
}

func GetNoteUrlByNotename(notename string) string {
	return UrlJoin(path.Join(GetNoteBase(), notename))
}

func GetTagBase() string {
	return "/" + conf.Conf.Site.TagUrl
}

func GetTagUrl(id int64) string {
	return UrlJoin(path.Join(GetTagBase(), strconv.FormatInt(id, 10)))
}

func GetTagNameUrl(name string) string {
	return UrlJoin(path.Join(GetTagBase(), name))
}

func GetPageBase() string {
	return "/" + conf.Conf.Site.PageUrl
}

func GetPageUrl(id int64) string {
	return UrlJoin(path.Join(GetPageBase(), strconv.FormatInt(id, 10)))
}

func GetRssBase() string {
	return "/" + conf.Conf.Site.RssUrl
}

func GetArchiveBase() string {
	return "/" + conf.Conf.Site.ArchiveUrl
}

func GetArchiveUrl() string {
	return UrlJoin(GetArchiveBase())
}

func GetArchiveMonthUrl(year, month int64) string {
	return UrlJoin(path.Join(GetArchiveBase(),
		strconv.FormatInt(year, 10)+"-"+strconv.FormatInt(month, 10)))
}

func GetSearchBase() string {
	return "/" + conf.Conf.Site.SearchUrl
}

func GetSearchUrl(keyword string) string {
	return UrlJoin(path.Join(GetSearchBase(), keyword))
}

func GetSearchPageUrl(keyword string, page int64) string {
	return UrlJoin(path.Join(GetSearchBase(), keyword, strconv.FormatInt(page, 10)))
}
