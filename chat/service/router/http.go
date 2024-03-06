package router

import (
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/fasthttp/websocket"
	"github.com/panjiawan/note/chat/conf"
	"go.uber.org/zap"
)

type HttpRouter struct {
	httpConf *conf.SysConf
	wsHandle *ihttp.Service
}

func New(cfg *conf.SysConf) *HttpRouter {
	return &HttpRouter{
		httpConf: cfg,
	}
}

// Run 启动函数
func (h *HttpRouter) Run() {
	// 启动WS
	h.wsHandle = ihttp.New(
		ihttp.WithAddress(h.httpConf.WsHost, h.httpConf.WsPort),
		ihttp.WithCertificate(h.httpConf.HttpsCertFile, h.httpConf.HttpsKeyFile),
		ihttp.WithRate(h.httpConf.RateLimitPerSec, h.httpConf.RrateLimitCapacity),
	)

	h.wsHandle.RegisterWS("/ws", websocket.TextMessage, cb)

	if err := h.wsHandle.Run(); err != nil {
		ilog.Error("ws server start error", zap.Error(err))
	}
}

func (h *HttpRouter) Close() {
}
