package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func InsertNoteTag(sr dbr.SessionRunner, note *model.Note, tag *model.Tag) error {
	noteTag := model.NewNoteTag(note.Id, tag.Id)
	_, err := sr.InsertInto(model.NoteTagTableName).
		Columns("note_id", "tag_id").Record(noteTag).Exec()
	return err
}

func DeleteNoteTagsNotExist(sr dbr.SessionRunner, note *model.Note, tagList []*model.Tag) error {
	if len(tagList) == 0 {
		return nil
	}
	tagIdList := make([]int64, 0, len(tagList))
	for _, tag := range tagList {
		tagIdList = append(tagIdList, tag.Id)
	}
	_, err := sr.DeleteFrom(model.NoteTagTableName).Where("note_id = ? and tag_id not in ?", note.Id, tagIdList).Exec()

	return err
}

func TagIdsByNoteId(sr dbr.SessionRunner, noteId int64) ([]int64, error) {
	return sr.Select("tag_id").From(model.NoteTagTableName).Where("note_id = ?", noteId).ReturnInt64s()
}

func NoteIdsByTagIds(sr dbr.SessionRunner, tagIds []int64) ([]int64, error) {
	return sr.Select("note_id").From(model.NoteTagTableName).Where("tag_id in ?", tagIds).ReturnInt64s()
}
