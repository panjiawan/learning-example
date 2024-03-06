package conf

import (
	"fmt"
	"github.com/coldwind/artist/pkg/icfg"
)

type SysConf struct {
	EnableDebug        bool   `yaml:"enableDebug"`
	EnableStdout       bool   `yaml:"enableStdout"`
	Https              bool   `yaml:"https"`
	HttpsCertFile      string `yaml:"httpsCertFile"`
	HttpsKeyFile       string `yaml:"httpsKeyFile"`
	HttpPort           int    `yaml:"httpPort"`
	GrpcPort           int    `yaml:"grpcPort"`
	WsHost             string `yaml:"wsHost"`
	WsPort             int    `yaml:"wsPort"`
	RateLimitPerSec    int    `yaml:"rateLimitPerSec"`
	RrateLimitCapacity int    `yaml:"rateLimitCapacity"`
	JwtSecret          string `yaml:"jwtSecret"`
	AvatarUrl          string `yaml:"avatarUrl"`
	RpcApiUrl          string `yaml:"rpcApiUrl"`
}

func (s *Handle) LoadSys() {
	path := fmt.Sprintf("%s/%s", s.path, "sys.yaml")
	err := icfg.Load(icfg.CfgTypeYaml, "sys", path, &SysConf{})
	if err != nil {
		panic(err)
	}
}

func (s *Handle) GetSysConf() *SysConf {
	return icfg.Get("sys").(*SysConf)
}
