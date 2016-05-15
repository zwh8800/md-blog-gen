package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func CheckIfNotenameDuplicated(tx *dbr.Tx, note *model.Note) (bool, error) {
	if _, err := tx.Select("unique_id").From(model.NoteTableName).
		Where("not unique_id = ? and notename = ?",
		note.UniqueId, note.Notename.String).ReturnInt64(); err != nil {
		if err != dbr.ErrNotFound {
			return false, err
		}
		return false, nil
	} else {
		return true, nil
	}
}

func InsertOrUpdateNote(tx *dbr.Tx, note *model.Note) error {
	dup, err := CheckIfNotenameDuplicated(tx, note)
	if err != nil {
		return err
	}
	if dup {
		note.Notename.String = ""
		note.Notename.Valid = false
	}

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
		note.Id = oldNote.Id
		if _, err := tx.Update(model.NoteTableName).Set("notename", note.Notename).
			Set("title", note.Title).Set("content", note.Content).
			Set("timestamp", note.Timestamp).Set("removed", note.Removed).
			Where("unique_id = ?", note.UniqueId).Exec(); err != nil {
			return err
		}
	} else {
		result, err := tx.InsertInto(model.NoteTableName).Columns("unique_id",
			"notename", "title", "url", "content", "timestamp").Record(note).Exec()
		if err != nil {
			return err
		}
		lastId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		note.Id = lastId
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

func NoteByNotename(sess *dbr.Session, notename string) (*model.Note, error) {
	note := &model.Note{}
	if err := sess.Select("*").From(model.NoteTableName).
		Where("notename = ? and removed is false", notename).LoadStruct(note); err != nil {
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

func NoteIdsByPage(sess *dbr.Session, page, limit int64) ([]int64, error) {
	offset := page * limit
	noteIdList := make([]int64, 0)
	if _, err := sess.Select("id").From(model.NoteTableName).
		Where("removed is false").OrderBy("timestamp desc").
		Offset(uint64(offset)).Limit(uint64(limit)).LoadValues(&noteIdList); err != nil {
		return nil, err
	}
	return noteIdList, nil
}

func NotesByTagId(sess *dbr.Session, tagId int64) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if _, err := sess.Select("Note.id", "Note.unique_id", "Note.notename", "Note.title",
		"Note.url", "Note.content", "Note.timestamp", "Note.removed").From(model.NoteTableName).
		Join(model.NoteTagTableName, "Note.id = NoteTag.note_id").Where("tag_id = ? and removed is false", tagId).
		LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func NotesByIds(sess *dbr.Session, noteIds []int64, limit int64) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if len(noteIds) == 0 {
		return noteList, nil
	}
	if _, err := sess.Select("*").From(model.NoteTableName).
		Where("id in ?", noteIds).OrderBy("timestamp desc").
		Limit(uint64(limit)).LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func NoteGroupByMonth(sess *dbr.Session) ([]*model.YearMonth, map[*model.YearMonth][]*model.Note, error) {
	monthList := make([]*model.YearMonth, 0)

	if _, err := sess.Select("YEAR(timestamp) year", "MONTH(timestamp) month").
		From(model.NoteTableName).Where("removed is false").
		GroupBy("YEAR(timestamp)", "MONTH(timestamp)").
		OrderBy("YEAR(timestamp), MONTH(timestamp) desc").LoadStructs(&monthList); err != nil {
		return nil, nil, err
	}

	noteListMap := make(map[*model.YearMonth][]*model.Note, len(monthList))
	for _, month := range monthList {
		noteList, err := NotesMonth(sess, month)
		if err != nil {
			return nil, nil, err

		}
		noteListMap[month] = noteList
	}

	return monthList, noteListMap, nil
}

func NotesMonth(sess *dbr.Session, month *model.YearMonth) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if _, err := sess.Select("*").From(model.NoteTableName).
		Where("removed is false and year(timestamp) = ? and month(timestamp) = ?", month.Year, month.Month).
		OrderBy("timestamp desc").LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}
