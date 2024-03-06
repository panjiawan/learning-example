package control

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type tokenEtl struct {
	Token string `json:"token"`
}

func GetToken(client *ihttp.WSClient, in gjson.Result) {
	uid := in.Get("uid").Int()
	if uid == 0 {
		client.Send(internal.PackError(code.ErrorGetToken))
		return
	}
	token := internal.GenToken(uint64(uid))

	// 发送成功消息
	client.Send(internal.PackOutput(code.OutTokenSuccess, &tokenEtl{
		Token: token,
	}))

	ilog.Info("get token", zap.Int64("connid", client.ConnId), zap.Uint64("uid", uint64(uid)), zap.String("token", token))
}
