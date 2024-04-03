package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjiawan/note/chat/conf"
	"time"

	"gitee.com/yqyn_live/idl/pb/api"
	"github.com/coldwind/artist/pkg/icache"
	"github.com/coldwind/artist/pkg/igrpc"
	"github.com/coldwind/artist/pkg/ilog"
	"go.uber.org/zap"
)

type Animation struct {
	Res string `json:"res"`
	Snd string `json:"snd"`
}

type UserDataCollect struct {
	Uid         uint64     `json:"uid"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Signature   string     `json:"signature"`
	Badges      []string   `json:"badges"`
	IsLocked    bool       `json:"isLocked"`
	IsManager   bool       `json:"isManager"`
	IsQuiet     bool       `json:"isQuiet"`
	QuietExpire int64      `json:"quiteExpire"`
	GuardAnim   *Animation `json:"guardAnim"`
	VehicleAnim *Animation `json:"vehicleAnim"`
	VehicleId   int32      `json:"vehicleId"`
	VehicleName string     `json:"vehicleName"`
}

func GetUserCollect(uid, roomId uint64) (*UserDataCollect, error) {
	lazyCache := icache.GetLazyCache("userdata", time.Second*10)
	key := fmt.Sprintf("%d_%d", uid, roomId)
	if reply, ok := lazyCache.Get(key); ok {
		return reply.(*UserDataCollect), nil
	}

	if conn, err := igrpc.GetGrpcClientHandle(conf.GetHandle().GetGrpcConf().Api); err == nil {
		defer conn.Close()
		req := &api.UserCollectRequest{
			Uid:    uid,
			RoomId: roomId,
		}
		apiService := api.NewApiClient(conn)
		if res, err := apiService.GetUserCollect(context.Background(), req); err == nil {
			if res.GetCodeMsg().GetCode() != 0 {
				ilog.Error("get collect error", zap.Any("errcode", res.GetCodeMsg()))
				return nil, errors.New(res.GetCodeMsg().GetMsg())
			}
			collectRes := &UserDataCollect{
				Uid:         res.GetInfo().GetUid(),
				Nickname:    res.GetInfo().GetNickname(),
				Avatar:      res.GetInfo().GetAvatar(),
				Badges:      res.GetInfo().GetBadges(),
				IsLocked:    res.GetInfo().GetIsLocked(),
				IsManager:   res.GetInfo().GetIsManager(),
				Signature:   res.GetInfo().GetSignature(),
				IsQuiet:     res.GetInfo().GetIsQuiet(),
				QuietExpire: res.GetInfo().GetQuietExpire(),
			}
			if res.GetInfo().GetGuardAnim() != nil {
				collectRes.GuardAnim = &Animation{
					Snd: res.GetInfo().GetGuardAnim().GetSnd(),
					Res: res.GetInfo().GetGuardAnim().GetRes(),
				}
			}
			if res.GetInfo().GetVehicleAnim() != nil {
				collectRes.VehicleAnim = &Animation{
					Snd: res.GetInfo().GetVehicleAnim().GetSnd(),
					Res: res.GetInfo().GetVehicleAnim().GetRes(),
				}
				collectRes.VehicleId = res.GetInfo().GetVehicleId()
				collectRes.VehicleName = res.GetInfo().GetVehicleName()
			}
			lazyCache.Set(key, collectRes, true)

			return collectRes, nil
		} else {
			ilog.Debug(" api grpc call error", zap.Any("err", err))
			return nil, err
		}
	} else {
		ilog.Debug("get api grpc handle error", zap.Any("err", err))
		return nil, err
	}
}
