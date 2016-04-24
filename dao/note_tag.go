package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func InsertNoteTag(tx *dbr.Tx, note *model.Note, tag *model.Tag) error {
	noteTag := model.NewNoteTag(note.Id, tag.Id)
	_, err := tx.InsertInto(model.NoteTagTableName).
		Columns("note_id", "tag_id").Record(noteTag).Exec()
	return err
}

func DeleteNoteTagsNotExist(tx *dbr.Tx, note *model.Note, tagList []*model.Tag) error {
	if len(tagList) == 0 {
		return nil
	}
	tagIdList := make([]int64, 0, len(tagList))
	for _, tag := range tagList {
		tagIdList = append(tagIdList, tag.Id)
	}
	_, err := tx.DeleteFrom(model.NoteTagTableName).Where("note_id = ? and tag_id not in ?", note.Id, tagIdList).Exec()

	return err
}

func TagIdsByNoteId(sess *dbr.Session, noteId int64) ([]int64, error) {
	return sess.Select("tag_id").From(model.NoteTagTableName).Where("note_id = ?", noteId).ReturnInt64s()
}

func NoteIdsByTagIds(sess *dbr.Session, tagIds []int64) ([]int64, error) {
	return sess.Select("note_id").From(model.NoteTagTableName).Where("tag_id in ?", tagIds).ReturnInt64s()
}
