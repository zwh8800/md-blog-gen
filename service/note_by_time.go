package service

import (
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
)

func NotesOrderByTime(page, limit int64) ([]*model.Note, map[int64][]*model.Tag, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NotesByPage(sess, page, limit)
	if err != nil {
		return nil, nil, 0, err
	}
	if len(noteList) == 0 {
		return noteList, nil, maxPage, nil
	}

	noteIdList := make([]int64, 0, len(noteList))
	for _, note := range noteList {
		noteIdList = append(noteIdList, note.Id)
	}
	tagListMap, err := dao.TagsByNoteIds(sess, noteIdList)
	if err != nil {
		return nil, nil, 0, err
	}
	return noteList, tagListMap, maxPage, nil
}

func NotesWithoutTagOrderByTime(page, limit int64) ([]*model.Note, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteList, err := dao.NotesByPage(sess, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return noteList, maxPage, nil
}

func NoteIdsOrderByTime(page, limit int64) ([]int64, int64, error) {
	sess := dbConn.NewSession(nil)
	noteCount, err := dao.CountNote(sess)
	if err != nil {
		return nil, 0, err
	}
	maxPage := (noteCount-1)/limit + 1

	page-- //数据库层的页数从0开始数
	noteIdList, err := dao.NoteIdsByPage(sess, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return noteIdList, maxPage, nil
}

func NoteGroupByMonth() ([]*model.YearMonth, map[*model.YearMonth][]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	monthList, err := dao.YearMonthList(sess, false)
	if err != nil {
		return nil, nil, err
	}

	noteListMap := make(map[*model.YearMonth][]*model.Note, len(monthList))
	for _, month := range monthList {
		noteList, err := dao.NotesMonth(sess, month)
		if err != nil {
			return nil, nil, err
		}
		noteListMap[month] = noteList
	}

	return monthList, noteListMap, nil
}

func findMonth(monthList []*model.YearMonth, item *model.YearMonth) int {
	for i := 0; i < len(monthList); i++ {
		if monthList[i].Year == item.Year &&
			monthList[i].Month == item.Month {
			return i
		}
	}
	return -1
}

func PrevNextMonth(month *model.YearMonth) (*model.YearMonth, *model.YearMonth, error) {
	sess := dbConn.NewSession(nil)
	monthList, err := dao.YearMonthList(sess, true)
	if err != nil {
		return nil, nil, err
	}
	i := findMonth(monthList, month)
	if i == -1 {
		return nil, nil, nil
	} else if i <= 0 && i >= len(monthList)-1 {
		return nil, nil, nil
	} else if i <= 0 {
		return nil, monthList[i+1], nil
	} else if i >= len(monthList)-1 {
		return monthList[i-1], nil, nil
	} else {
		return monthList[i-1], monthList[i+1], nil
	}
}

func NotesByMonth(month *model.YearMonth) ([]*model.Note, error) {
	sess := dbConn.NewSession(nil)
	return dao.NotesMonth(sess, month)
}
