package util

import (
	"net/url"
	"path"
	"strconv"

	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
)

func GetNoteBase() string {
	return "/" + conf.Conf.Site.NoteUrl
}

func GetNoteUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetNoteBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetNoteUrlByNotename(notename string) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetNoteBase(), notename)
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetTagBase() string {
	return "/" + conf.Conf.Site.TagUrl
}

func GetTagUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetTagBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetTagNameUrl(name string) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetTagBase(), name)
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetPageBase() string {
	return "/" + conf.Conf.Site.PageUrl
}

func GetPageUrl(id int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetPageBase(), strconv.FormatInt(id, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetRssBase() string {
	return "/" + conf.Conf.Site.RssUrl
}

func GetArchiveBase() string {
	return "/" + conf.Conf.Site.ArchiveUrl
}

func GetArchiveUrl() string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := GetArchiveBase()
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetArchiveMonthUrl(month *model.YearMonth) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetArchiveBase(),
		strconv.FormatInt(month.Year, 10)+"-"+strconv.FormatInt(month.Month, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetSearchBase() string {
	return "/" + conf.Conf.Site.SearchUrl
}

func GetSearchUrl(keyword string) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetSearchBase(), keyword)
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}

func GetSearchPageUrl(keyword string, page int64) string {
	baseUrl, err := url.Parse(conf.Conf.Site.BaseUrl)
	if err != nil {
		glog.Error(err)
		return ""
	}
	noteChild := path.Join(GetSearchBase(), keyword, strconv.FormatInt(page, 10))
	u, err := url.Parse(noteChild)
	if err != nil {
		glog.Error(err)
		return ""
	}
	return baseUrl.ResolveReference(u).String()
}
