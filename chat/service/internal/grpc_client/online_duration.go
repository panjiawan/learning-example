package grpc_client

import (
	"context"
	"gitee.com/yqyn_live/idl/pb/api"
	"github.com/coldwind/artist/pkg/igrpc"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/conf"
	"go.uber.org/zap"
)

func OnlineDuration(uid, roomId uint64) error {
	ilog.Debug("OnlineDuration--------------------------", zap.Uint64("uid", uid), zap.Uint64("roomId", roomId))
	if conn, err := igrpc.GetGrpcClientHandle(conf.GetHandle().GetGrpcConf().Api); err == nil {
		defer conn.Close()
		apiService := api.NewApiClient(conn)
		_, err := apiService.OnlineDuration(context.Background(), &api.OnlineDurationRequest{
			Uid:    uid,
			RoomId: roomId,
		})
		if err != nil {
			ilog.Error("OnlineDuration error", zap.Any("err", err))
			return err
		}

		return nil

	} else {
		ilog.Debug("get api grpc handle error", zap.Any("err", err))
		return err
	}
}
