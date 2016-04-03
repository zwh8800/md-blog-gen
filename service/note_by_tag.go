package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesByTagId(id int64) (*model.Tag, []*model.Note, error) {
	sess := dbConn.NewSession(nil)
	tag, err := dao.TagById(sess, id)
	if err != nil {
		return nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, err
	}
	return tag, noteList, nil
}

func NotesByTagName(name string) (*model.Tag, []*model.Note, error) {
	sess := dbConn.NewSession(nil)
	tag, err := dao.TagByName(sess, name)
	if err != nil {
		return nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, err
	}
	return tag, noteList, nil
}

func AllNotesTags() ([]*model.Tag, map[int64][]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	tagList, err := dao.Tags(sess)
	if err != nil {
		return nil, nil, err
	}
	noteListMap := make(map[int64][]*model.Note)
	for i, tag := range tagList {
		noteList, err := dao.NotesByTagId(sess, tag.Id)
		if err != nil {
			tagList = append(tagList[:i], tagList[i+1:]...)
			continue
		} else if len(noteList) == 0 {
			tagList = append(tagList[:i], tagList[i+1:]...)
			continue
		}
		noteListMap[tag.Id] = noteList
	}

	return tagList, noteListMap, nil
}
