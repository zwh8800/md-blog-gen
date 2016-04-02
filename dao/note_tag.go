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
	_, err := tx.DeleteFrom(model.NoteTagTableName).Where("tag_id not in ?", tagIdList).Exec()

	return err
}
