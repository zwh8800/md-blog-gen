package spider

import (
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
)

func findBlogTagUl(doc *goquery.Document) *goquery.Selection {
	var result *goquery.Selection
	doc.Find("#file-list").Children().Each(func(i int, s *goquery.Selection) {
		ul := s.Find("ul")
		tagName, _ := ul.Find(".tag-item.item").Attr("tag-name")
		if tagName == conf.Conf.Spider.SpiderTag {
			result = ul
		}
	})
	return result
}

func findNoteList(doc *goquery.Document) []*model.Note {
	ul := findBlogTagUl(doc)
	liList := ul.Find(".file-item.item")
	result := make([]*model.Note, 0, liList.Length())
	liList.Each(func(i int, s *goquery.Selection) {
		timestampStr, _ := s.Attr("file-created-date")
		timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			glog.Warning(err)
			return
		}
		a := s.Find("a")
		url, _ := a.Attr("href")
		span := a.Find("span")
		idStr, _ := span.Attr("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			glog.Warning(err)
			return
		}
		title := span.Text()
		result = append(result, model.NewNote(id, title, url, timestamp))
	})
	return result
}

func findTagListMap(doc *goquery.Document) map[int64][]*model.Tag {
	tagListMap := make(map[int64][]*model.Tag)
	doc.Find("#file-list .tag-list").Each(func(i int, s *goquery.Selection) {
		tagSel := s.Find(".tag-item.item")
		notesSel := s.Find(".file-item.item")
		tagName, _ := tagSel.Attr("tag-name")
		tag := model.NewTag(tagName)
		if tagName == conf.Conf.Spider.SpiderTag {
			return
		}

		notesSel.Each(func(i int, s *goquery.Selection) {
			span := s.Find("a span")
			idStr, _ := span.Attr("id")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				glog.Warning(err)
				return
			}
			tagList, ok := tagListMap[id]
			if !ok {
				tagList = make([]*model.Tag, 0)
			}
			tagList = append(tagList, tag)
			tagListMap[id] = tagList
		})

	})

	return tagListMap
}

func findStartPage() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(conf.Conf.Spider.StartUrl)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
