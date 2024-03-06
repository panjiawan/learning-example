package router

import (
	"github.com/panjiawan/note/chat/service/code"
	"github.com/valyala/fasthttp"
)

func (h *HttpRouter) Filter(ctx *fasthttp.RequestCtx) code.OutputCode {
	return code.Success
}
