package code

type logicCode struct {
	Code int
	Msg  string
}

type OutputCode *logicCode

var Success OutputCode = &logicCode{
	Code: 0,
	Msg:  "",
}
