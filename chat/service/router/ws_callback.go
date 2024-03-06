package router

import (
	"errors"
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/fasthttp/websocket"
	"github.com/panjiawan/note/chat/service/control"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/control/mgr"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type WsCallback struct {
}

var cb = &WsCallback{}

func (w *WsCallback) OnConnect(ctx *fasthttp.RequestCtx, client *ihttp.WSClient) error {
	if client == nil {
		return errors.New("client connect error:client is nil")
	}
	ilog.Debug("on connect", zap.Int64("connid", client.ConnId))

	// 初始化user data数据
	client.UserData = &common.UserData{}
	mgr.UserMgr.Add(client)

	return nil
}

func (w *WsCallback) OnMessage(client *ihttp.WSClient, msgType int, msg []byte) {
	if msgType != websocket.TextMessage {
		return
	}
	control.MessageOperate(client, msg)
}

func (w *WsCallback) OnClose(client *ihttp.WSClient) {
	ilog.Info("offline", zap.Int64("connid", client.ConnId))

	if client.UserData != nil {
		userData := client.UserData.(*common.UserData)
		if userData.RoomId > 0 && userData.Uid > 0 {
			//	room.LeaveRoom(userData.RoomId, userData.Uid, client.ConnId)
			userData.RoomId = 0
		}
	}
	mgr.UserMgr.Del(client)
}
