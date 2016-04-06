package service

import (
	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesByTagId(id int64) (*model.Tag, []*model.Note, map[int64][]*model.Tag, error) {
	sess := dbConn.NewSession(nil)
	tag, err := dao.TagById(sess, id)
	if err != nil {
		return nil, nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, nil, err
	}
	tagListMap := make(map[int64][]*model.Tag)
	for _, note := range noteList {
		tags, err := dao.TagsByNoteId(sess, note.Id)
		if err != nil {
			glog.Warning(err)
			tags = make([]*model.Tag, 0)
		}
		tagListMap[note.Id] = tags
	}
	return tag, noteList, tagListMap, nil
}

func NotesByTagName(name string) (*model.Tag, []*model.Note, map[int64][]*model.Tag, error) {
	sess := dbConn.NewSession(nil)
	tag, err := dao.TagByName(sess, name)
	if err != nil {
		return nil, nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, nil, err
	}
	tagListMap := make(map[int64][]*model.Tag)
	for _, note := range noteList {
		tags, err := dao.TagsByNoteId(sess, note.Id)
		if err != nil {
			glog.Warning(err)
			tags = make([]*model.Tag, 0)
		}
		tagListMap[note.Id] = tags
	}
	return tag, noteList, tagListMap, nil
}

func AllNotesTags() ([]*model.Tag, map[int64][]*model.Note, map[int64][]*model.Tag, error) {
	sess := dbConn.NewSession(nil)
	tagList, err := dao.Tags(sess)
	if err != nil {
		return nil, nil, nil, err
	}
	noteListMap := make(map[int64][]*model.Note)
	tagListMap := make(map[int64][]*model.Tag)
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
		for _, note := range noteList {
			if _, exists := tagListMap[note.Id]; exists {
				continue
			}

			tags, err := dao.TagsByNoteId(sess, note.Id)
			if err != nil {
				glog.Warning(err)
				tags = make([]*model.Tag, 0)
			}
			tagListMap[note.Id] = tags
		}
	}

	return tagList, noteListMap, tagListMap, nil
}
