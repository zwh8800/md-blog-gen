package spider

import (
	"runtime/debug"

	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/service"
)

func Go() {
	glog.Info("spider started")
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln("panic in spider recovered:", err, string(debug.Stack()))
		}
	}()

	startPage, err := findStartPage()
	if err != nil {
		glog.Error(err)
		return
	}
	noteList := findNoteList(startPage)
	tagListMap := findTagListMap(startPage)
	for _, note := range noteList {
		FindNoteContent(note)
	}

	if err := service.SaveNoteList(noteList, tagListMap); err != nil {
		glog.Error(err)
		return
	}

	glog.Info("spider finished")
}
