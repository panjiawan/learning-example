package room

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/coldwind/artist/pkg/iutils"
	"github.com/panjiawan/note/chat/service/code"
	"github.com/panjiawan/note/chat/service/control/common"
	"github.com/panjiawan/note/chat/service/defines"
	"github.com/panjiawan/note/chat/service/dispatch"
	"github.com/panjiawan/note/chat/service/internal"
	"github.com/panjiawan/note/chat/service/internal/grpc_client"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

type etlOutExitRoom struct {
	Uid uint64 `json:"uid"`
	Num int    `json:"num"`
}

type etlOutJoinRoom struct {
	Uid         int64                  `json:"uid"`
	Nickname    string                 `json:"nickname"`
	Avatar      string                 `json:"avatar"`
	Badges      []string               `json:"badges"`
	IsLocked    bool                   `json:"isLocked"`
	IsManager   bool                   `json:"isManager"`
	Play        *grpc_client.Animation `json:"play"`
	Num         int                    `json:"num"`
	VehicleName string                 `json:"vehicleName"`
}

type AnnouncementEtl struct {
	List []*AnnouncementItem `json:"list"`
}
type AnnouncementItem struct {
	Content string `json:"content"`
	Color   string `json:"color"`
}

func (r *RoomManage) JoinRoom(client *ihttp.WSClient, in gjson.Result) {
	// 获取房间ID号
	roomId, err := iutils.StringToUint64(in.Get("roomId").String())
	if err != nil || roomId == 0 {
		client.Send(internal.PackError(code.ErrorRoomParam))
		return
	}

	userData := client.UserData.(*common.UserData)
	fromRoomId := userData.RoomId
	uid := userData.Uid
	connId := client.ConnId
	roomModel := &model.RoomModel{}

	// 获取用户信息集合
	userCollectData, err := grpc_client.GetUserCollect(userData.Uid, roomId)
	if err == nil && len(userCollectData.Badges) > 0 {
		userData.Badges = userCollectData.Badges
	}

	if err != nil {
		client.Send(internal.PackError(code.ErrorUserNotExist))
		return
	}

	// 清除上一个直播间用户的连接
	if fromRoomId > 0 {
		fromRoom := r.Get(fromRoomId)
		if fromRoom != nil {
			fromRoom.Lock()
			if _, ok := fromRoom.Audience[uid]; ok {
				delete(fromRoom.Audience[uid], connId)
				ilog.Debug("delete conn from old room", zap.Uint64("uid", uid), zap.Int64("connid", connId))
				if len(fromRoom.Audience[uid]) == 0 {
					delete(fromRoom.Audience, uid)

					// 删除用户房间redis记录
					roomModel.Offline(fromRoomId, uid)

					// 发送给老直播间用户下线通知
					onlineNum, _ := roomModel.CountOnlineNum(fromRoomId)
					res := &etlOutExitRoom{
						Uid: uid,
						Num: onlineNum,
					}
					dispatch.PublishToRoom(fromRoomId, internal.PackOutput(code.OutUserExit, res))
					ilog.Debug("delete uid from old room", zap.Uint64("uid", uid), zap.Int64("connid", connId))
				}
			}
			fromRoom.Unlock()
		}
	}

	// 添加到新的直播间里
	room := r.Get(roomId)
	if room == nil {
		client.Send(internal.PackError(code.ErrorRoomNotExist))
		return
	}

	if userCollectData.IsQuiet && userCollectData.QuietExpire > time.Now().Unix() {
		r.Quiet(roomId, uid, userCollectData.QuietExpire)
	} else {
		if r.IsQuiet(roomId, uid) {
			r.CancelQuiet(roomId, uid)
		}
	}

	room.Lock()
	defer room.Unlock()
	if _, ok := room.Audience[uid]; !ok {
		room.Audience[uid] = map[int64]*ihttp.WSClient{
			client.ConnId: client,
		}
	} else {
		room.Audience[uid][client.ConnId] = client
	}

	// 更新用户数据
	userData.RefreshTime = time.Now().Unix()
	userData.RoomId = roomId

	if roomId != userData.Uid {
		// 设置上线redis存放数据
		roomModel.Online(roomId, uid)

		// 发送用户进入通知&用户人数变动数
		userModel := &model.UserModel{}
		user, err := userModel.Get(userData.Uid)
		if err != nil {
			return
		}

		onlineNum, _ := roomModel.CountOnlineNum(roomId)
		res := &etlOutJoinRoom{
			Uid:       int64(user.ID),
			Nickname:  user.Nickname,
			Avatar:    defines.AbsoluteUrl(user.Avatar),
			Num:       onlineNum,
			Badges:    userData.Badges,
			IsLocked:  userCollectData.IsLocked,
			IsManager: userCollectData.IsManager,
		}

		ilog.Debug("userCollectData", zap.Any("userCollectData", userCollectData))
		if userCollectData.GuardAnim != nil {
			// 播放守护
			res.Play = &grpc_client.Animation{
				Res: defines.AbsoluteUrl(userCollectData.GuardAnim.Res),
				Snd: defines.AbsoluteUrl(userCollectData.GuardAnim.Snd),
			}
		} else if userCollectData.VehicleAnim != nil {
			// 坐骑
			res.Play = &grpc_client.Animation{
				Res: defines.AbsoluteUrl(userCollectData.VehicleAnim.Res),
				Snd: defines.AbsoluteUrl(userCollectData.VehicleAnim.Snd),
			}
			res.VehicleName = userCollectData.VehicleName
		}

		// 发送给本人公告信息
		announceModel := &model.AnnouncementModel{}
		list := announceModel.GetList()
		if len(list) > 0 {
			announceList := make([]*AnnouncementItem, 0, len(list))
			for _, v := range list {
				announceList = append(announceList, &AnnouncementItem{
					Content: v.Content,
					Color:   v.Color,
				})
			}
			client.Send(internal.PackOutput(code.OutAnnouncement, &AnnouncementEtl{
				List: announceList,
			}))
		}

		// 发送用户数变动数据
		dispatch.PublishToRoom(roomId, internal.PackOutput(code.OutUserJoin, res))
		ilog.Debug("join room success", zap.Uint64("roomId", roomId), zap.Uint64("uid", user.ID))
	}
}
