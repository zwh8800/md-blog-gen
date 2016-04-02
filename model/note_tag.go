package model

const NoteTagTableName = "NoteTag"

type NoteTag struct {
	NoteId int64 `db:"note_id" json:"noteId" struct:"note_id"`
	TagId  int64 `db:"tag_id" json:"tagId" struct:"tag_id"`
}

func NewNoteTag(noteId int64, tagId int64) *NoteTag {
	return &NoteTag{
		NoteId: noteId,
		TagId:  tagId,
	}
}
