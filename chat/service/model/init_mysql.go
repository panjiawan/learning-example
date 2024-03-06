package model

import (
	"github.com/coldwind/artist/pkg/imysql"
	"github.com/panjiawan/note/chat/conf"
	"github.com/panjiawan/workaholic/pkg/plog"
	"go.uber.org/zap"
)

func initMysql(mysqlConf *conf.MysqlConf) {
	mysqlHandles = make(map[string]*imysql.Service)

	for _, cfg := range mysqlConf.Hosts {
		mysqlHandles[cfg.Name] = imysql.New(
			imysql.WithConnection(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DB),
			imysql.WithLimit(cfg.MaxIdle, cfg.MaxOpen),
			imysql.WithPrefix(cfg.Prefix),
			imysql.WithDebug(cfg.Debug),
			imysql.WithCharset("utf8mb4"),
		)
	}

	for k, f := range mysqlHandles {
		if err := f.Run(); err != nil {
			plog.Error("mysql start error", zap.String("key", k), zap.Error(err))
			panic(err)
		} else {
			plog.Info("mysql started", zap.String("key", k))
		}
	}
}
