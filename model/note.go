package model

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/util"
)

const NoteTableName = "Note"

type Note struct {
	Id           int64          `db:"id" json:"id" struct:"id"`
	UniqueId     int64          `db:"unique_id" json:"uniqueId" struct:"unique_id"`
	Notename     dbr.NullString `db:"notename" json:"notename" struct:"notename"`
	Title        string         `db:"title" json:"title" struct:"title"`
	Url          string         `db:"url" json:"url" struct:"url"`
	Content      string         `db:"content" json:"content" struct:"content"`
	Timestamp    time.Time      `db:"timestamp" json:"timestamp" struct:"timestamp"`
	LastModified time.Time      `db:"last_modified" json:"lastModified" struct:"lastModified"`
	Hash         string         `db:"hash" json:"hash" struct:"hash"`
	Removed      bool           `db:"removed" json:"removed" struct:"removed"`
}

func NewNote(uniqueId int64, title string, url string, timestamp time.Time) *Note {
	return &Note{
		UniqueId:     uniqueId,
		Title:        title,
		Url:          url,
		Timestamp:    timestamp,
		LastModified: time.Now().UTC(),
		Removed:      false,
	}
}

func (obj *Note) UnescapedContent() template.HTML {
	return template.HTML(obj.Content)
}

func (obj *Note) FillContent(content string) {
	obj.Content = content
	obj.Hash = util.MD5(content)
}

func (obj *Note) Preview() string {
	src := obj.Content
	buffer := bytes.NewBufferString(src)

	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return ""
	}

	return goquery.NewDocumentFromNode(doc.Find("p").Get(1)).Text()
}

func (obj *Note) FormattedTimestamp() string {
	return obj.Timestamp.Local().Format("2006-01-02 15:04 PM")
}

func (obj *Note) FormattedDate() string {
	return obj.Timestamp.Local().Format("Jan 02, 2006")
}

func (obj *Note) FormattedLastModified() string {
	return obj.LastModified.Local().Format("2006-01-02 15:04:05 PM")
}

func (obj *Note) YearMonth() string {
	return obj.Timestamp.Local().Format("2006-1")
}

func (obj *Note) PreviewWithKeyword(keyword string) template.HTML {
	preview := obj.Preview()
	return template.HTML(strings.Replace(preview,
		keyword, `<span class="red">`+keyword+`</span>`, -1))
}

func (obj *Note) TitleWithKeyword(keyword string) template.HTML {
	return template.HTML(strings.Replace(obj.Title,
		keyword, `<span class="red">`+keyword+`</span>`, -1))
}
