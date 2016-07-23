package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func IsNoteModified(sr dbr.SessionRunner, note *model.Note) (bool, error) {
	hash, err := sr.Select("hash").From(model.NoteTableName).
		Where("unique_id = ?", note.UniqueId).ReturnString()
	if err != nil {
		if err == dbr.ErrNotFound {
			return true, nil
		}
		return false, err
	}
	return hash != note.Hash, nil
}

func CheckIfNotenameDuplicated(sr dbr.SessionRunner, note *model.Note) (bool, error) {
	if _, err := sr.Select("unique_id").From(model.NoteTableName).
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

func InsertOrUpdateNote(sr dbr.SessionRunner, note *model.Note) error {
	dup, err := CheckIfNotenameDuplicated(sr, note)
	if err != nil {
		return err
	}
	if dup {
		note.Notename.String = ""
		note.Notename.Valid = false
	}

	oldNote := &model.Note{}
	found := true
	if err := sr.Select("*").From(model.NoteTableName).Where("unique_id = ?", note.UniqueId).
		LoadStruct(oldNote); err != nil {
		if err == dbr.ErrNotFound {
			found = false
		} else {
			return err
		}
	}
	if found {
		note.Id = oldNote.Id
		if _, err := sr.Update(model.NoteTableName).Set("notename", note.Notename).
			Set("title", note.Title).Set("content", note.Content).
			Set("timestamp", note.Timestamp).Set("last_modified", note.LastModified).
			Set("hash", note.Hash).Set("removed", note.Removed).
			Where("unique_id = ?", note.UniqueId).Exec(); err != nil {
			return err
		}
	} else {
		result, err := sr.InsertInto(model.NoteTableName).Columns("unique_id",
			"notename", "title", "url", "content", "timestamp", "last_modified", "hash").Record(note).Exec()
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

func NoteById(sr dbr.SessionRunner, id int64) (*model.Note, error) {
	note := &model.Note{}
	if err := sr.Select("*").From(model.NoteTableName).
		Where("id = ? and removed is false", id).LoadStruct(note); err != nil {
		return nil, err
	}

	return note, nil
}

func NoteByNotename(sr dbr.SessionRunner, notename string) (*model.Note, error) {
	note := &model.Note{}
	if err := sr.Select("*").From(model.NoteTableName).
		Where("notename = ? and removed is false", notename).LoadStruct(note); err != nil {
		return nil, err
	}

	return note, nil
}

func RemoveUnpublishedNote(sr dbr.SessionRunner, noteList []*model.Note) error {
	uniqueIdList := make([]int64, 0, len(noteList))
	for _, note := range noteList {
		uniqueIdList = append(uniqueIdList, note.UniqueId)
	}
	if _, err := sr.Update(model.NoteTableName).Set("removed", true).
		Where("removed is false and unique_id not in ?", uniqueIdList).Exec(); err != nil {
		return err
	}
	return nil
}

func CountNote(sr dbr.SessionRunner) (int64, error) {
	return sr.Select("count(*)").From(model.NoteTableName).
		Where("removed is false").ReturnInt64()
}

func NotesByPage(sr dbr.SessionRunner, page, limit int64) ([]*model.Note, error) {
	offset := page * limit
	noteList := make([]*model.Note, 0)
	if _, err := sr.Select("*").From(model.NoteTableName).
		Where("removed is false").OrderBy("timestamp desc").
		Offset(uint64(offset)).Limit(uint64(limit)).LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func NoteIdsByPage(sr dbr.SessionRunner, page, limit int64) ([]int64, error) {
	offset := page * limit
	noteIdList := make([]int64, 0)
	if _, err := sr.Select("id").From(model.NoteTableName).
		Where("removed is false").OrderBy("timestamp desc").
		Offset(uint64(offset)).Limit(uint64(limit)).LoadValues(&noteIdList); err != nil {
		return nil, err
	}
	return noteIdList, nil
}

func NotesByTagId(sr dbr.SessionRunner, tagId int64) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if _, err := sr.Select("Note.id", "Note.unique_id", "Note.notename", "Note.title",
		"Note.url", "Note.content", "Note.timestamp", "Note.removed").From(model.NoteTableName).
		Join(model.NoteTagTableName, "Note.id = NoteTag.note_id").Where("tag_id = ? and removed is false", tagId).
		LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func NotesByIds(sr dbr.SessionRunner, noteIds []int64, limit int64) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if len(noteIds) == 0 {
		return noteList, nil
	}
	if _, err := sr.Select("*").From(model.NoteTableName).
		Where("id in ?", noteIds).OrderBy("timestamp desc").
		Limit(uint64(limit)).LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func YearMonthList(sr dbr.SessionRunner, isAsc bool) ([]*model.YearMonth, error) {
	monthList := make([]*model.YearMonth, 0)

	if _, err := sr.Select("YEAR(timestamp) year", "MONTH(timestamp) month").
		From(model.NoteTableName).Where("removed is false").
		GroupBy("YEAR(timestamp)", "MONTH(timestamp)").
		OrderDir("YEAR(timestamp), MONTH(timestamp)", isAsc).LoadStructs(&monthList); err != nil {
		return nil, err
	}
	return monthList, nil
}

func NotesMonth(sr dbr.SessionRunner, month *model.YearMonth) ([]*model.Note, error) {
	noteList := make([]*model.Note, 0)
	if _, err := sr.Select("*").From(model.NoteTableName).
		Where("removed is false and year(timestamp) = ? and month(timestamp) = ?", month.Year, month.Month).
		OrderBy("timestamp desc").LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}

func CountNoteByKeyword(sr dbr.SessionRunner, keyword string) (int64, error) {
	search := "%" + keyword + "%"
	return sr.Select("count(*)").From(model.NoteTableName).
		Where("removed is false and (title like ? or notename like ? or content like ?)", search, search, search).
		ReturnInt64()
}

func NoteByKeyword(sr dbr.SessionRunner, keyword string, page, limit int64) ([]*model.Note, error) {
	offset := page * limit
	search := "%" + keyword + "%"

	noteList := make([]*model.Note, 0)
	if _, err := sr.Select("*").From(model.NoteTableName).
		Where("removed is false and (title like ? or notename like ? or content like ?)", search, search, search).
		OrderBy("timestamp desc").Offset(uint64(offset)).Limit(uint64(limit)).
		LoadStructs(&noteList); err != nil {
		return nil, err
	}
	return noteList, nil
}
