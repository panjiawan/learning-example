package common

type UserData struct {
	Uid         uint64
	RoomId      uint64
	Nickname    string
	Avatar      string
	Level       int
	Badges      []string
	RefreshTime int64
}

func (u *UserData) IsAuth() bool {
	return u.Uid > 0
}
