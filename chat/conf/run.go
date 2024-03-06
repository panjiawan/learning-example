package conf

type Handle struct {
	path string
}

var handler *Handle = nil

func New(etcPath string) *Handle {
	if handler == nil {
		handler = &Handle{
			path: etcPath,
		}
	}

	return handler
}

func GetHandle() *Handle {
	return handler
}

func (s *Handle) Run() {
	s.LoadSys()
	s.LoadMysql()
	s.LoadRedis()
}

func (s *Handle) Close() {
}
