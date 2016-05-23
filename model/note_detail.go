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
