package boot

import (
	"github.com/coldwind/artist/pkg/isignal"
	"github.com/panjiawan/workaholic/pkg/plog"
	"go.uber.org/zap"
	"os"
)

func closeSignalListen() {
	defer func() {
		if err := recover(); err != nil {
			plog.Error("signal listen error", zap.Any("err", err))
		}
	}()

	signal := isignal.New()
	signal.Register(os.Interrupt, func(signal os.Signal, args interface{}) {
		close()
		os.Exit(0)
	})
	signal.Listen()
}
