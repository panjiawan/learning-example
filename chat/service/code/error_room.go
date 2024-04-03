package code

var (
	ErrorRoomParam OutputCode = &logicCode{
		Code: 10001,
		Msg:  "房间参数有误",
	}
	ErrorRoomNotExist OutputCode = &logicCode{
		Code: 10002,
		Msg:  "房间号不存在",
	}
	ErrorMsg OutputCode = &logicCode{
		Code: 10003,
		Msg:  "消息内容不合法",
	}
	ErrorQuiet OutputCode = &logicCode{
		Code: 10004,
		Msg:  "你已经禁言，不能发消息",
	}
	ErrorLevelLimit OutputCode = &logicCode{
		Code: 10006,
		Msg:  "未达到聊天的等级",
	}
)
