package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func TagsByNoteId(sess *dbr.Session, noteId int64) ([]*model.Tag, error) {
	tagList := make([]*model.Tag, 0)
	if _, err := sess.Select("id", "name").From(model.TagTableName).
		Join(model.NoteTagTableName, "Tag.id = NoteTag.tag_id").Where("note_id = ?", noteId).
		LoadStructs(&tagList); err != nil {
		return nil, err
	}
	return tagList, nil
}

func SelectTagOrInsertIfNotExists(tx *dbr.Tx, tag *model.Tag) (*model.Tag, error) {
	found := true
	if err := tx.Select("*").From(model.TagTableName).
		Where("name = ?", tag.Name).LoadStruct(tag); err != nil {
		if err == dbr.ErrNotFound {
			found = false
		} else {
			return nil, err
		}
	}
	if !found {
		result, err := tx.InsertInto(model.TagTableName).Columns("name").Record(tag).Exec()
		if err != nil {
			return nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		tag.Id = id
	}

	return tag, nil
}
