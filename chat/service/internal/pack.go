package internal

import (
	"github.com/coldwind/artist/pkg/iutils"
	"github.com/panjiawan/note/chat/service/code"
)

// 输出到前端的数据结构
type ErrorMsg struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type OutputMsg struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

var DefaultErrorCode int32 = 400

func PackOutput(cmd string, data interface{}) []byte {
	out := &OutputMsg{
		Cmd:  cmd,
		Data: data,
	}

	return iutils.JsonMarshal(out)
}

func PackError(err code.OutputCode) []byte {
	out := &ErrorMsg{
		Msg:  err.Msg,
		Code: err.Code,
	}

	return iutils.JsonMarshal(out)
}
