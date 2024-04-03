package model

import (
	"context"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ISyncMsg interface {
	Watch(func(msg string))
	Publish(msg string) error
}

func SyncMsgNew() ISyncMsg {
	return &SyncMsgWithRedis{
		channel: "wnwdkj:sync:msg",
	}
}

type SyncMsgWithRedis struct {
	channel string
}

func (s *SyncMsgWithRedis) Watch(fn func(msg string)) {
	ilog.Info("SyncMsgWithRedis start")
	ctx := context.Background()
	sub := redisHandles[redisSyncMsgKey].GetConn().Subscribe(ctx, s.channel)
	defer func() {
		sub.Unsubscribe(ctx, s.channel)
		sub.Close()
	}()

	_, err := sub.Receive(context.Background())
	if err != nil {
		panic(err)
	}

	// 检测收到的消息类型
	for msg := range sub.Channel(redis.WithChannelSize(4096)) {
		// ilog.Debug("sync recv", zap.String("msg", msg.Payload))
		fn(msg.Payload)
	}
	ilog.Error("SyncMsgWithRedis error")

	go s.Watch(fn)
}

func (s *SyncMsgWithRedis) Publish(msg string) error {
	err := redisHandles[redisSyncMsgKey].GetConn().Publish(context.TODO(), s.channel, msg).Err()
	ilog.Debug("publish call redis", zap.String("channel", s.channel), zap.Error(err))
	if err != nil {
		ilog.Error("publish call redis error", zap.String("channel", s.channel), zap.Error(err))
	}
	return err
}
