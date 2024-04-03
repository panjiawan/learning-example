package room

import (
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/coldwind/artist/pkg/iutils"
	"github.com/panjiawan/note/chat/service/control/mgr"
	"github.com/panjiawan/note/chat/service/dispatch"
	"go.uber.org/zap"
	"time"
)

func Run() {
	dispatch.RegisterCallback(dispatch.DCmdRoom, Dispatcher)
	dispatch.RegisterCallback(dispatch.DCmdRoomByUid, DispatcherByUid)

	dispatch.RegisterCallback(dispatch.SCmdQuiet, Quiet)
	dispatch.RegisterCallback(dispatch.SCmdCancelQuiet, CancelQuiet)
}

// 推送消息到指定直播间
func Dispatcher(data *dispatch.Msg) {
	if data.RoomId > 0 && data.Uid == 0 {
		RoomHandle.SendToRoom(data.RoomId, data.Data)
	}
}

// 推送消息到用户所在的直播间内
func DispatcherByUid(data *dispatch.Msg) {
	if data.Uid > 0 {
		roomIds := mgr.UserMgr.GetRoomIdsByUid(data.Uid)
		for _, roomId := range roomIds {
			RoomHandle.SendToRoom(roomId, data.Data)
		}
	}
}

// 禁言调用
func Quiet(data *dispatch.Msg) {
	ilog.Debug("Quiet callback", zap.Any("data", data))
	if data.RoomId > 0 && data.Uid > 0 {
		// 解析数据
		expire := iutils.BytesToInt64(data.Data)
		if expire > time.Now().Unix() {
			RoomHandle.Quiet(data.RoomId, data.Uid, expire)
		}
	}
}

func CancelQuiet(data *dispatch.Msg) {
	ilog.Debug("CancelQuiet callback", zap.Any("data", data))
	if data.RoomId > 0 && data.Uid > 0 {
		RoomHandle.CancelQuiet(data.RoomId, data.Uid)
	}
}
