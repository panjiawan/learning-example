package room

import (
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/dispatch"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/panjiawan/note/chat/service/model"
	"go.uber.org/zap"
)

type etlLeaveRoom struct {
	Uid uint64 `json:"uid"`
	Num int    `json:"num"`
}

func LeaveRoom(roomId, uid uint64, connId int64) {
	if isExist, err := RoomHandle.RemoveUser(roomId, uid, connId); err == nil {
		if isExist {
			// 发送离开房间信息
			roomModel := &model.RoomModel{}
			onlineNum, _ := roomModel.CountOnlineNum(roomId)
			res := &etlLeaveRoom{
				Uid: uid,
				Num: onlineNum,
			}
			dispatch.PublishToRoom(roomId, internal.PackOutput(code.OutUserExit, res))
			ilog.Debug("leave room success", zap.Uint64("roomId", roomId), zap.Uint64("uid", uid))
		}
	}
}
