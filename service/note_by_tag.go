package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesByTagId(id int64) (*model.Tag, []*model.Note, map[int64][]*model.Tag, error) {
	sess := newSession()
	tag, err := dao.TagById(sess, id)
	if err != nil {
		return nil, nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	noteIdList := make([]int64, 0)
	for _, note := range noteList {
		noteIdList = append(noteIdList, note.Id)
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, nil, err
	}
	return tag, noteList, tagListMap, nil
}

func NotesByTagName(name string) (*model.Tag, []*model.Note, map[int64][]*model.Tag, error) {
	sess := newSession()
	tag, err := dao.TagByName(sess, name)
	if err != nil {
		return nil, nil, nil, err
	}
	noteList, err := dao.NotesByTagId(sess, tag.Id)
	if err != nil {
		return nil, nil, nil, err
	}

	noteIdList := make([]int64, 0)
	for _, note := range noteList {
		noteIdList = append(noteIdList, note.Id)
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, nil, err
	}
	return tag, noteList, tagListMap, nil
}

func AllNotesTags() ([]*model.Tag, map[int64][]*model.Note, map[int64][]*model.Tag, error) {
	sess := newSession()
	tagList, err := dao.Tags(sess)
	if err != nil {
		return nil, nil, nil, err
	}

	noteIdList := make([]int64, 0)
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
		for _, note := range noteList {
			noteIdList = append(noteIdList, note.Id)
		}
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, nil, err
	}

	return tagList, noteListMap, tagListMap, nil
}
