package conf

import (
	"fmt"
	"github.com/coldwind/artist/pkg/icfg"
)

type GrpcConf struct {
	Api []string `yaml:"api"`
	Ws  []string `yaml:"ws"`
	Pk  []string `yaml:"pk"`
}

func (s *Handle) LoadGrpc() {
	path := fmt.Sprintf("%s/%s", s.path, "grpc.yaml")
	err := icfg.Load(icfg.CfgTypeYaml, "grpc", path, &GrpcConf{})
	if err != nil {
		panic(err)
	}
}

func (s *Handle) GetGrpcConf() *GrpcConf {
	return icfg.Get("grpc").(*GrpcConf)
}
