package defines

import (
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/conf"
	"go.uber.org/zap"
)

func AbsoluteUrl(url string) string {
	if len(url) <= 7 {
		if url == "" {
			return url
		}
		return conf.GetHandle().GetSysConf().AvatarUrl + url
	}

	protocol := url[:7]
	ilog.Debug("protocol", zap.String("protocol", protocol), zap.String("AvatarUrl", conf.GetHandle().GetSysConf().AvatarUrl))

	if protocol == "http://" || protocol == "https:/" {
		return url
	}

	return conf.GetHandle().GetSysConf().AvatarUrl + url
}
