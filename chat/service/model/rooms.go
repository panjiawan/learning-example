package model

import (
	"context"
	"fmt"
	"github.com/coldwind/artist/pkg/iredis"
	"github.com/redis/go-redis/v9"
	"time"
)

type Room struct {
	ID          uint
	Name        string `gorm:"column:name"`
	IsLiving    int    `gorm:"column:is_living"`
	IsPk        int    `gorm:"column:is_pk"`
	Locked      int    `gorm:"column:locked"`
	PkStatus    int    `gorm:"column:pk_status"`
	PkId        int    `gorm:"column:pk_id"`
	PkUserId    int    `gorm:"column:pk_user_id"`
	Level       int    `gorm:"column:level"`
	Exp         int    `gorm:"column:exp"`
	AudienceNum int    `gorm:"column:audience_num"`
	GuardNum    int    `gorm:"column:guard_num"`
	ManagerNum  int    `gorm:"column:manager_num"`
	Sort        int    `gorm:"column:sort"`
	Status      int    `gorm:"column:status"` // 状态 1为审核通过，0为待审核
}

type RoomModel struct {
}

var onlineKeyFormat = "online:%d"

func (Room) TableName() string {
	return "rooms"
}

func (a *RoomModel) Get(id uint64) (*Room, error) {
	res := &Room{}
	db := mysqlHandles[mysqlLiveKey].Handle().Where("id=?", id).First(res)
	if db.Error != nil {
		return nil, db.Error
	}

	return res, nil
}

// 用于在线用户统计的代码
func (a *RoomModel) Online(roomId, uid uint64) error {
	key := getKey(fmt.Sprintf(onlineKeyFormat, roomId))
	_ = redisHandles[redisLiveKey].GetConn().ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: uid,
	})
	return nil
}

// 删除在线用户
func (a *RoomModel) Offline(roomId, uid uint64) error {
	key := getKey(fmt.Sprintf(onlineKeyFormat, roomId))
	redisHandles[redisLiveKey].GetConn().ZRem(context.Background(), key, uid)

	return nil
}

// 用于在线用户统计的代码
func (a *RoomModel) GetOnlineUid(roomId uint64, begin, limit int) ([]uint64, error) {
	endTime := time.Now().Unix()
	startTime := endTime - 60

	key := getKey(fmt.Sprintf(onlineKeyFormat, roomId))
	sliceCmd := redisHandles[redisLiveKey].GetConn().ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", startTime),
		Max:    fmt.Sprintf("%d", endTime),
		Offset: int64(begin),
		Count:  int64(limit),
	})
	onlines := make([]uint64, 0, limit)
	iredis.New()
	if res, err := sliceCmd.Result(); err == nil {
		onlines = iredis.StringSliceToUInt64Slice(res)
		return onlines, err
	} else {
		return onlines, err
	}
}

// 用于在线用户统计的代码
func (a *RoomModel) CountOnlineNum(roomId uint64) (int, error) {
	endTime := time.Now().Unix()
	startTime := endTime - 60
	key := getKey(fmt.Sprintf(onlineKeyFormat, roomId))
	num, err := redisHandles[redisLiveKey].GetConn().ZCount(context.Background(), key, fmt.Sprintf("%d", startTime), fmt.Sprintf("%d", endTime)).Result()
	if err != nil {
		return 0, err
	}

	return int(num), err
}
