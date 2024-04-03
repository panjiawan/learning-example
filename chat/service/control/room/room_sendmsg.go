package room

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/defines"
	"github.com/panjiawan/note/chat/service/dispatch"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

type etlSendMsg struct {
	Uid      uint64   `json:"uid"`
	Nickname string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Msg      string   `json:"msg"`
	Badges   []string `json:"badges"`
}

var msgLen = 200

func (r *RoomManage) SendMsg(client *ihttp.WSClient, in gjson.Result) {
	// 获取消息内容
	msg := in.Get("msg").String()
	if msg == "" || len(msg) > msgLen {
		client.Send(internal.PackError(code.ErrorMsg))
		return
	}

	userData := client.UserData.(*common.UserData)
	// 未达到等级&不是主播自己直播间
	if userData.Level < levelLimit && userData.Uid != userData.RoomId {
		client.Send(internal.PackError(code.ErrorLevelLimit))
		return
	}
	userData.RefreshTime = time.Now().Unix()

	room := r.Get(userData.RoomId)
	if room == nil {
		client.Send(internal.PackError(code.ErrorRoomNotExist))
		return
	}

	// 判断禁言
	if r.IsQuiet(userData.RoomId, userData.Uid) {
		client.Send(internal.PackError(code.ErrorQuiet))
		return
	}

	// 发送
	res := &etlSendMsg{
		Uid:      userData.Uid,
		Nickname: userData.Nickname,
		Avatar:   defines.AbsoluteUrl(userData.Avatar),
		Msg:      msg,
		Badges:   userData.Badges,
	}

	// 发送用户数变动数据
	dispatch.PublishToRoom(userData.RoomId, internal.PackOutput(code.OutRecvMsg, res))
	ilog.Debug("send message success", zap.Uint64("roomId", userData.RoomId), zap.Uint64("uid", userData.Uid))
}
