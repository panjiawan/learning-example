package room

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/tidwall/gjson"
)

func (r *RoomManage) ExitRoom(client *ihttp.WSClient, in gjson.Result) {
	// 获取房间ID号
	userData := client.UserData.(*common.UserData)
	roomId := userData.RoomId
	uid := userData.Uid
	connId := client.ConnId

	if uid == 0 || connId == 0 || roomId == 0 {
		client.Send(internal.PackError(code.ErrorRoomParam))
		return
	}

	LeaveRoom(roomId, uid, connId)
}
