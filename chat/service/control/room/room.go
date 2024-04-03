package room

import (
	"errors"
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/model"
	"go.uber.org/zap"
	"sync"
)

type Room struct {
	sync.RWMutex
	RoomId     uint64                               // 房间ID
	RoomNo     string                               // 直播间房号 -- 目前与RoomId值一致类型不同
	CompereUid uint64                               // 主持人
	Audience   map[uint64]map[int64]*ihttp.WSClient // 观众
	Quiet      sync.Map                             // 禁言列表
}

type RoomManage struct {
	rooms sync.Map
}

var RoomHandle = &RoomManage{}

// get room info
func (r *RoomManage) Get(roomId uint64) *Room {

	if room, ok := r.rooms.Load(roomId); ok {
		return room.(*Room)
	}

	// 读取room数据
	roomModel := &model.RoomModel{}
	roomDB, err := roomModel.Get(roomId)
	if err != nil || roomDB == nil {
		return nil
	}

	room := &Room{
		RoomId:     roomId,
		CompereUid: roomId,
		Audience:   make(map[uint64]map[int64]*ihttp.WSClient),
	}
	r.rooms.Store(roomId, room)

	return room
}

// delete
// when bool is true, user exit room
func (r *RoomManage) RemoveUser(roomId uint64, uid uint64, connId int64) (bool, error) {
	room := r.Get(roomId)
	isExist := false
	if room == nil {
		return isExist, errors.New("room not exist")
	}

	room.Lock()
	defer room.Unlock()
	if _, ok := room.Audience[uid]; !ok {
		return true, nil
	}

	if _, ok := room.Audience[uid][connId]; !ok {
		return isExist, nil
	}
	delete(room.Audience[uid], connId)
	ilog.Debug("remove connid from room", zap.Uint64("roomId", roomId), zap.Int64("connid", connId))

	if len(room.Audience[uid]) == 0 {
		delete(room.Audience, uid)
		isExist = true
		// 删除redis数据
		roomModel := &model.RoomModel{}
		roomModel.Offline(roomId, uid)
		ilog.Debug("remove uid from room", zap.Uint64("roomId", roomId), zap.Uint64("uid", uid))
	}

	return isExist, nil
}

// 发送给房间内的所有用户
func (r *RoomManage) SendToRoom(roomId uint64, msg []byte) {
	if value, ok := r.rooms.Load(roomId); ok {
		room := value.(*Room)
		room.RLock()
		defer room.RUnlock()
		for _, users := range room.Audience {
			for _, conn := range users {
				if conn.UserData.(*common.UserData).RoomId == roomId {
					conn.Send(msg)
				}
			}
		}
	}
}

// 广播给所有房间的人
func (r *RoomManage) Broadcast(msg []byte) {
	r.rooms.Range(func(key, value interface{}) bool {
		room := value.(*Room)
		room.RLock()
		defer room.RUnlock()
		for _, users := range room.Audience {
			for _, conn := range users {
				conn.Send(msg)
			}
		}
		return true
	})
}
