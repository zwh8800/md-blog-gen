package dao

import (
	"github.com/gocraft/dbr"

	"github.com/zwh8800/md-blog-gen/model"
)

func TagById(sess *dbr.Session, id int64) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := sess.Select("*").From(model.TagTableName).
		Where("id = ?", id).LoadStruct(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func TagByName(sess *dbr.Session, name string) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := sess.Select("*").From(model.TagTableName).
		Where("name = ?", name).LoadStruct(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func Tags(sess *dbr.Session) ([]*model.Tag, error) {
	tagList := make([]*model.Tag, 0)
	if _, err := sess.Select("*").From(model.TagTableName).
		LoadStructs(&tagList); err != nil {
		return nil, err
	}
	return tagList, nil
}

func TagsByNoteId(sess *dbr.Session, noteId int64) ([]*model.Tag, error) {
	tagList := make([]*model.Tag, 0)
	if _, err := sess.Select("Tag.id", "Tag.name").From(model.TagTableName).
		Join(model.NoteTagTableName, "Tag.id = NoteTag.tag_id").Where("note_id = ?", noteId).
		LoadStructs(&tagList); err != nil {
		return nil, err
	}
	return tagList, nil
}

func TagsByNoteIds(sess *dbr.Session, noteIds []int64) (map[int64][]*model.Tag, error) {
	noteIdTagList := make([]*model.NoteIdTag, 0)

	if _, err := sess.Select("NoteTag.note_id note_id", "Tag.id tag_id", "Tag.Name tag_name").
		From(model.TagTableName).Join(model.NoteTagTableName, "Tag.id = NoteTag.tag_id").
		Where("note_id in ?", noteIds).LoadStructs(&noteIdTagList); err != nil {
		return nil, err
	}

	tagListMap := make(map[int64][]*model.Tag)
	for _, noteIdTag := range noteIdTagList {
		tagList, ok := tagListMap[noteIdTag.NoteId]
		if !ok {
			tagList = make([]*model.Tag, 0)
		}
		tagList = append(tagList, &model.Tag{
			Id:   noteIdTag.TagId,
			Name: noteIdTag.TagName,
		})
		tagListMap[noteIdTag.NoteId] = tagList
	}

	return tagListMap, nil
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
