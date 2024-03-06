package common

import (
	"sync"
	"testing"
)

func TestUserDataClear(t *testing.T) {
	var b sync.Map
	u := &UserData{
		Uid:    10,
		RoomId: 5,
	}
	b.Store("user", u)
	u.Uid = 2

	e, _ := b.Load("user")
	t.Log(e.(*UserData).Uid)
}
