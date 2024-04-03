package mgr

import "github.com/panjiawan/note/chat/service/dispatch"

func Run() {
	dispatch.RegisterCallback(dispatch.DCmdMgrSend, DispatchToUser)
	dispatch.RegisterCallback(dispatch.DCmdBroadcast, DispatchToAll)
	dispatch.RegisterCallback(dispatch.DCmdBroadcastWithoutRoom, DispatchToAllWithoutRoom)
}

func DispatchToUser(data *dispatch.Msg) {
	if data.Uid > 0 && len(data.Data) > 0 {
		UserMgr.SendToUser(data.Uid, data.Data)
	}
}

func DispatchToAll(data *dispatch.Msg) {
	if len(data.Data) > 0 {
		UserMgr.Broadcast(data.Data)
	}
}

func DispatchToAllWithoutRoom(data *dispatch.Msg) {
	if len(data.Data) > 0 {
		UserMgr.BroadcastWithoutRoom(data.RoomId, data.Data)
	}
}
