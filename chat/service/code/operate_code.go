package code

var (
	GetToken    = "getToken"
	InOnline    = "online"
	InJoinRoom  = "joinRoom"
	InSendMsg   = "sendMsg"
	InLogin     = "login" // 登录
	InLeaveRoom = "leaveRoom"
)

var (
	OutTokenSuccess = "tokenSuccess" // 获取token成功
	OutLoginSuccess = "loginSuccess" // WS登录成功
	OutUserJoin     = "userJoin"     // 有新用户加入房间
	OutUserExit     = "userExit"     // 有用户退出房间
	OutRecvMsg      = "recvMsg"      // 接收到消息
	OutAnnouncement = "announcement" // 公告
)
