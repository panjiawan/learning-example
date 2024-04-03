package room

import (
	"fmt"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/coldwind/artist/pkg/iutils"
	"github.com/panjiawan/note/chat/service/model"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

var levelLimit = 0

func LevelLimitDaemon() {
	defer func() {
		if e := recover(); e != nil {
			ilog.Error("level limit error", zap.String("stack", string(debug.Stack())))
		}
	}()

	for {
		// 获取等级限制
		settingModel := model.SettingModel{}
		if info, err := settingModel.Get(); err == nil {
			baseRes := ParseContent(info.Base)
			if chatLevel, ok := baseRes["chat_level"]; ok {
				if level, err := iutils.StringToInt(chatLevel); err == nil {
					levelLimit = level
				}
			}
		}
		time.Sleep(time.Second * 30)
	}
}

func ParseContent(content string) map[string]string {
	config := make(map[string]string)
	gRes := gjson.Get(content, "#")
	for i := 0; i < int(gRes.Int()); i++ {
		keyName := fmt.Sprintf("%d.key", i)
		valueName := fmt.Sprintf("%d.value", i)
		cfgKey := gjson.Get(content, keyName).String()
		cfgValue := gjson.Get(content, valueName).String()
		if cfgKey == "" {
			continue
		}
		config[cfgKey] = cfgValue
	}

	return config
}
