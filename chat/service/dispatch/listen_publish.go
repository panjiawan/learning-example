package dispatch

import (
	"encoding/json"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/model"
	"go.uber.org/zap"
)

func (d *Dispatch) Listen(fn func(*Msg)) {
	// 坚听分布式消息内容
	d.handler = model.SyncMsgNew()
	d.handler.Watch(func(msg string) {
		res := &Msg{}
		json.Unmarshal([]byte(msg), res)
		fn(res)
	})
}

func (d *Dispatch) Publish(res *Msg) {
	msgByte, _ := json.Marshal(res)
	ilog.Debug("publish", zap.String("msg", string(msgByte)))

	if err := d.handler.Publish(string(msgByte)); err != nil {
		ilog.Error("publish error", zap.Error(err))
	}
}
