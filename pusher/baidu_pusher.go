package pusher

import (
	"math"

	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func BaiduPush() {
	glog.Infoln("Baidu push start")
	noteList, _, _, err := service.NotesOrderByTime(0, math.MaxInt64)
	if err != nil {
		glog.Error(err)
		return
	}
	urls := make([]string, 0, len(noteList))
	for _, note := range noteList {
		urls = append(urls, util.GetNoteUrl(note.Id))
	}
	if err := util.PushUrlToBaidu(urls); err != nil {
		glog.Error(err)
		return
	}
	glog.Infoln("Baidu push finish")
}
