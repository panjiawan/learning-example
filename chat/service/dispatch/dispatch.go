package dispatch

import "github.com/panjiawan/note/chat/service/model"

type Callback func(*Msg)

type Dispatch struct {
	handler model.ISyncMsg
	fn      map[string]Callback
}

type Msg struct {
	Cmd    string
	RoomId uint64
	Uid    uint64
	Data   []byte
}

var dispatcher *Dispatch

var (
	DCmdRoom                 = "room"
	DCmdBroadcast            = "broadcast"
	DCmdBroadcastWithoutRoom = "broadcastWithRoom"
	DCmdRoomByUid            = "roomByUid"

	DCmdMgrSend = "mgrSend"

	SCmdQuiet       = "quiet"
	SCmdCancelQuiet = "cancelQuiet"
)

func Run() {
	dispatcher = &Dispatch{
		fn: map[string]Callback{},
	}
	go dispatcher.Listen(ReceiveOperate)
}

// 收到消息后处理回调
func ReceiveOperate(data *Msg) {
	if fn, ok := dispatcher.fn[data.Cmd]; ok {
		fn(data)
	}
}

// 注册回调
func RegisterCallback(cmd string, cb Callback) {
	dispatcher.fn[cmd] = cb
}

// 往房间推送消息
func PublishToRoom(roomId uint64, data []byte) {
	dispatcher.Publish(&Msg{
		Cmd:    DCmdRoom,
		RoomId: roomId,
		Data:   data,
	})
}

// 往用户推送消息
func PublishToUser(uid uint64, data []byte) {
	dispatcher.Publish(&Msg{
		Cmd:  DCmdMgrSend,
		Uid:  uid,
		Data: data,
	})
}

// 往用户推送消息
func Broadcast(data []byte) {
	dispatcher.Publish(&Msg{
		Cmd:  DCmdBroadcast,
		Data: data,
	})
}

// 剔除指定room后往用户推送消息
func BroadcastWithoutRoom(roomId uint64, data []byte) {
	dispatcher.Publish(&Msg{
		Cmd:    DCmdBroadcastWithoutRoom,
		RoomId: roomId,
		Data:   data,
	})
}

// 往用户所在的房间推送消息
func PublishToRoomByUid(uid uint64, data []byte) {
	dispatcher.Publish(&Msg{
		Cmd:  DCmdRoomByUid,
		Uid:  uid,
		Data: data,
	})
}
