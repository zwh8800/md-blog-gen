package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func Tags() ([]*model.Tag, error) {
	sess := newSession()
	tagList, err := dao.Tags(sess)
	if err != nil {
		return nil, err
	}

	return tagList, nil
}

func TagsByCount() ([]*model.Tag, error) {
	sess := newSession()
	tagList, err := dao.TagsByCount(sess)
	if err != nil {
		return nil, err
	}

	return tagList, nil
}

func TagsByNoteId(noteId int64) ([]*model.Tag, error) {
	sess := newSession()
	return dao.TagsByNoteId(sess, noteId)
}
