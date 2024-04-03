package room

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/coldwind/artist/pkg/iutils"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/panjiawan/note/chat/service/internal/grpc_client"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

// 在线状态调整
func (r *RoomManage) Online(client *ihttp.WSClient, in gjson.Result) {
	// 读取用户在房间中的数据
	userData := client.UserData.(*common.UserData)
	if userData.Uid == 0 {
		return
	}
	// 刷新在线时间
	userData.RefreshTime = time.Now().Unix()

	// 获取房间ID号
	roomId, err := iutils.StringToUint64(in.Get("roomId").String())
	if err != nil || roomId == 0 {
		return
	}

	room := r.Get(roomId)
	if room == nil {
		client.Send(internal.PackError(code.ErrorRoomNotExist))
		return
	}

	room.RLock()
	defer room.RUnlock()
	if _, ok := room.Audience[userData.Uid]; !ok {
		return
	}

	if _, ok := room.Audience[userData.Uid][client.ConnId]; !ok {
		return
	}

	if len(room.Audience[userData.Uid]) == 0 {
		delete(room.Audience, userData.Uid)
		return
	}
	ilog.Debug("room online", zap.Uint64("roomId", roomId), zap.Uint64("uid", userData.Uid))

	if roomId != userData.Uid {
		roomModel := &model.RoomModel{}
		roomModel.Online(roomId, userData.Uid)
	} else {
		// 主播则同步在线状态 api会进行统计
		grpc_client.OnlineDuration(userData.Uid, roomId)
	}
}
