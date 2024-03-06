package model

import (
	"github.com/coldwind/artist/pkg/imysql"
	"github.com/coldwind/artist/pkg/iredis"
	"github.com/panjiawan/note/chat/conf"
)

var (
	mysqlHandles map[string]*imysql.Service
	redisHandles map[string]*iredis.Service

	mysqlLiveKey = "live"

	redisPrefix     = "live_database"
	redisLiveKey    = "live"
	redisSyncMsgKey = "sync_msg"
)

func Run(mysqlCfg *conf.MysqlConf, redisCfg *conf.RedisConf) {
	initMysql(mysqlCfg)
	initRedis(redisCfg)
}
