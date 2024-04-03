package model

import (
	"github.com/coldwind/artist/pkg/icache"
	"github.com/coldwind/artist/pkg/ilog"
	"go.uber.org/zap"
	"time"
)

var announcementConfigCache *icache.LazyCacheItem

func init() {
	announcementConfigCache = icache.GetLazyCache("announcementConfig", time.Minute*5)
}

type Announcement struct {
	ID      uint64
	Content string `gorm:"column:content"`
	Color   string `gorm:"column:color"`
}

type AnnouncementModel struct {
}

func (a *AnnouncementModel) GetList() []*Announcement {
	if res, ok := announcementConfigCache.Get("config"); ok {
		return res.([]*Announcement)
	}

	res := make([]*Announcement, 0, 10)
	db := mysqlHandles["live"].Handle().Order("id desc").Find(&res)
	if db.Error != nil {
		ilog.Error("Announcement", zap.Error(db.Error))
	}

	announcementConfigCache.Set("config", res, true)

	return res
}
