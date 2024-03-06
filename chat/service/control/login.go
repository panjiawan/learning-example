package control

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/tidwall/gjson"
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
}
