package controller

import (
	"errors"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/zwh8800/md-blog-gen/controller/index"
	"github.com/zwh8800/md-blog-gen/controller/middleware/cache"
	"github.com/zwh8800/md-blog-gen/controller/opensearch"
	"github.com/zwh8800/md-blog-gen/controller/rss"
	"github.com/zwh8800/md-blog-gen/controller/sitemap"
	"github.com/zwh8800/md-blog-gen/util"
)

func Route(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		index.ErrorHandler(c, 404, errors.New("404 Not Found"))
	})

	indexGroup := r.Group("/", cache.Middleware())
	{
		indexGroup.GET("/", index.Index)
		indexGroup.GET(path.Join(util.GetPageBase(), ":page"), index.Index)
		indexGroup.GET(util.GetTagBase(), index.AllTag)
		indexGroup.GET(path.Join(util.GetTagBase(), ":tag"), index.Tag)
		indexGroup.GET(path.Join(util.GetNoteBase(), ":id"), index.Note)
		indexGroup.GET(util.GetArchiveBase(), index.Archive)
		indexGroup.GET(path.Join(util.GetArchiveBase(), ":month"), index.ArchiveMonth)
		indexGroup.GET(util.GetSearchBase(), index.SearchIndex)
		indexGroup.GET(path.Join(util.GetSearchBase(), ":keyword"), index.Search)
		indexGroup.GET(path.Join(util.GetSearchBase(), ":keyword", ":page"), index.Search)
		indexGroup.GET(path.Join("api", util.GetSearchBase(), ":keyword"), index.SearchTitle)
	}

	r.GET("/search.xml", opensearch.OpenSearch)

	r.GET(util.GetRssBase(), rss.Rss)
	r.GET("/.rss", rss.Rss)
	r.GET("/feed", rss.Rss)
	r.GET("/atom", rss.Atom)

	r.GET("/sitemap.xml", sitemap.SiteMap)

	r.Static("/static", "./static")
}
