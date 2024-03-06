package model

import (
	"fmt"
	"github.com/coldwind/artist/pkg/iredis"
	"github.com/panjiawan/note/chat/conf"
	"github.com/panjiawan/workaholic/pkg/plog"
	"go.uber.org/zap"
	"time"
)

func initRedis(redisConf *conf.RedisConf) {
	redisHandles = make(map[string]*iredis.Service)

	for _, cfg := range redisConf.Hosts {
		timeout := time.Duration(cfg.Timeout) * time.Second
		if cfg.Name == redisSyncMsgKey {
			timeout = 0
		}
		redisHandles[cfg.Name] = iredis.New(
			iredis.WithConnection(cfg.Host, cfg.Port),
			iredis.WithAuth(cfg.Auth),
			iredis.WithLimit(cfg.MinIdle, cfg.MaxIdle),
			iredis.WithReadTimeout(timeout),
			iredis.WithWriteTimeout(timeout),
		)
	}

	for k, f := range redisHandles {
		if err := f.Run(); err != nil {
			plog.Error("redis start error", zap.String("key", k), zap.Error(err))
			panic(err)
		} else {
			plog.Info("redis started", zap.String("key", k))
		}
	}
}

func getKey(key string) string {
	return fmt.Sprintf("%s:%s", redisPrefix, key)
}
