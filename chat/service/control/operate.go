package control

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/tidwall/gjson"
)

func MessageOperate(client *ihttp.WSClient, msg []byte) {
	operate := gjson.Get(string(msg), "cmd")
	if operate.String() == "Null" {
		return
	}
	in := gjson.Get(string(msg), "param")

	switch operate.String() {
	case code.InLogin:
		// 登录
		Login(client, in)
		////case code.InOnline:
		////	// 房间在线心跳包
		////	room.RoomHandle.Online(client, in)
		////case code.InJoinRoom:
		////	// 用户加入房间
		////	room.RoomHandle.JoinRoom(client, in)
		//case code.InSendMsg:
		//	// 房间内发送消息
		//	room.RoomHandle.SendMsg(client, in)
		//case code.InLeaveRoom:
		//	// 离开房间
		//	room.RoomHandle.ExitRoom(client, in)
	}
}
