package model

type NoteIdTag struct {
	TagId   int64  `db:"tag_id"`
	TagName string `db:"tag_name"`
	NoteId  int64  `db:"note_id"`
}
