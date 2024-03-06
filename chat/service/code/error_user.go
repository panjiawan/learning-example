package code

var (
	ErrorUserNotExist OutputCode = &logicCode{
		Code: 20001,
		Msg:  "用户信息不存在",
	}

	ErrorLoginFailure OutputCode = &logicCode{
		Code: 20002,
		Msg:  "用户登录失败",
	}

	ErrorToken OutputCode = &logicCode{
		Code: 20003,
		Msg:  "token验证失败",
	}

	ErrorGetToken OutputCode = &logicCode{
		Code: 20004,
		Msg:  "生成token失败",
	}
)
