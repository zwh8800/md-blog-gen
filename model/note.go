package model

import (
	"bytes"
	"html/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocraft/dbr"
)

const NoteTableName = "Note"

type Note struct {
	Id        int64          `db:"id" json:"id" struct:"id"`
	UniqueId  int64          `db:"unique_id" json:"uniqueId" struct:"unique_id"`
	Notename  dbr.NullString `db:"notename" json:"notename" struct:"notename"`
	Title     string         `db:"title" json:"title" struct:"title"`
	Url       string         `db:"url" json:"url" struct:"url"`
	Content   string         `db:"content" json:"content" struct:"content"`
	Timestamp time.Time      `db:"timestamp" json:"timestamp" struct:"timestamp"`
	Removed   bool           `db:"removed" json:"removed" struct:"removed"`
}

func NewNote(uniqueId int64, title string, url string, timestamp time.Time) *Note {
	return &Note{
		UniqueId:  uniqueId,
		Title:     title,
		Url:       url,
		Timestamp: timestamp,
		Removed:   false,
	}
}

func (obj *Note) UnescapedContent() template.HTML {
	return template.HTML(obj.Content)
}

func (obj *Note) FillContent(content string) {
	obj.Content = content
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

func (obj *Note) YearMonth() string {
	return obj.Timestamp.Local().Format("2006-1")
}
