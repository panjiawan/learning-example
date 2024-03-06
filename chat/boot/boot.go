package boot

import (
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/conf"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/panjiawan/note/chat/service/router"
	"github.com/panjiawan/workaholic/pkg/plog"
)

func Start(etcPath string, logPath string) {
	// load conf
	confHandle := conf.New(etcPath)
	confHandle.Run()

	// start log
	plog.Start(logPath, "note_chat.log", confHandle.GetSysConf().EnableDebug, confHandle.GetSysConf().EnableStdout)
	ilog.Start(logPath, "note_chat.log", confHandle.GetSysConf().EnableDebug, confHandle.GetSysConf().EnableStdout)

	plog.Info("conf started")

	// start signal
	go closeSignalListen()

	// 启动model
	model.Run(confHandle.GetMysqlConf(), confHandle.GetRedisConf())
	plog.Info("model started")

	// 启动消息分发坚听
	//dispatch.Run()

	// 启动控制器注册
	//control.Run()

	// 启动grpc
	//grpc_server.Run()

	route := router.New(confHandle.GetSysConf())
	route.Run()
}

func Stay() {
	select {}
}

// 优雅关闭调用点
func close() {

}
