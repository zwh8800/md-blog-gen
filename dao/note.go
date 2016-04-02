package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func InsertOrUpdateNote(tx *dbr.Tx, note *model.Note) error {
	oldNote := &model.Note{}
	found := true
	if err := tx.Select("*").From(model.NoteTableName).Where("unique_id = ?", note.UniqueId).
		LoadStruct(oldNote); err != nil {
		if err == dbr.ErrNotFound {
			found = false
		} else {
			return err
		}
	}
	if found {
		if _, err := tx.Update(model.NoteTableName).
			Set("title", note.Title).Set("content", note.Content).
			Set("timestamp", note.Timestamp).Set("removed", note.Removed).
			Where("unique_id = ?", note.UniqueId).Exec(); err != nil {
			return err
		}
	} else {
		if _, err := tx.InsertInto(model.NoteTableName).Columns("unique_id",
			"title", "url", "content", "timestamp").Record(note).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func NoteById(sess *dbr.Session, id int64) (*model.Note, error) {
	note := &model.Note{}
	if err := sess.Select("*").From(model.NoteTableName).
		Where("id = ? and removed is false", id).LoadStruct(note); err != nil {
		return nil, err
	}

	return note, nil
}

func RemoveUnpublishedNote(tx *dbr.Tx, noteList []*model.Note) error {
	uniqueIdList := make([]int64, 0, len(noteList))
	for _, note := range noteList {
		uniqueIdList = append(uniqueIdList, note.UniqueId)
	}
	if _, err := tx.Update(model.NoteTableName).Set("removed", true).
		Where("removed is false and unique_id not in ?", uniqueIdList).Exec(); err != nil {
		return err
	}
	return nil
}

func CountNote(sess *dbr.Session) (int64, error) {
	return sess.Select("count(*)").From(model.NoteTableName).ReturnInt64()
}

func NotesByPage(sess *dbr.Session, page, limit int64) ([]*model.Note, error) {
	offset := page * limit
	noteList := make([]*model.Note, 0)
	if _, err := sess.Select("*").From(model.NoteTableName).
		Where("removed is false").OrderBy("timestamp desc").
		Offset(uint64(offset)).Limit(uint64(limit)).LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}
