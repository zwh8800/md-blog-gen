package model

import (
	"time"

	"github.com/gocraft/dbr"
)

type NoteDetail struct {
	Id           int64          `json:"id"`
	Notename     dbr.NullString `json:"notename"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Timestamp    time.Time      `json:"timestamp"`
	LastModified time.Time      `json:"lastModified"`
	TagList      []string       `json:"tagList"`
}

func (obj *NoteDetail) FormattedTimestamp() string {
	return obj.Timestamp.Local().Format("2006-01-02 15:04 PM")
}

func (obj *NoteDetail) YearMonth() string {
	return obj.Timestamp.Local().Format("2006-1")
}
