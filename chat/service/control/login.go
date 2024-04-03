package control

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/control/mgr"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

type LoginEtl struct {
	Uid uint64 `json:"uid"`
}

func Login(client *ihttp.WSClient, in gjson.Result) {
	token := in.Get("token").String()
	if token == "" {
		client.Send(internal.PackError(code.ErrorToken))
		return
	}

	tokenInfo, err := internal.ParseToken(token)
	if err != nil || tokenInfo.Uid == 0 {
		client.Send(internal.PackError(code.ErrorToken))
		return
	}

	uid := tokenInfo.Uid
	if uid == 0 {
		client.Send(internal.PackError(code.ErrorToken))
		return
	}

	userData := &common.UserData{
		Uid:         uid,
		RoomId:      0,
		RefreshTime: time.Now().Unix(),
	}

	// 读用户数据
	userModel := &model.UserModel{}
	userInfo, err := userModel.Get(uid)
	if err != nil {
		client.Send(internal.PackError(code.ErrorUserNotExist))
		return
	}

	userData.Nickname = userInfo.Nickname
	userData.Avatar = userInfo.Avatar
	userData.Level = userInfo.Level
	client.UserData = userData

	//认证登陆
	if err = mgr.UserMgr.Authed(client); err != nil {
		client.Send(internal.PackError(code.ErrorLoginFailure))
		return
	}

	// 发送登录成功消息
	client.Send(internal.PackOutput(code.OutLoginSuccess, &LoginEtl{
		Uid: userData.Uid,
	}))
	ilog.Info("user login", zap.Int64("connid", client.ConnId), zap.Uint64("uid", uid), zap.String("nickname", userInfo.Nickname))
}
