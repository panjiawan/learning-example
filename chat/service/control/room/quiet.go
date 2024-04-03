package room

import "time"

// 禁言
func (r *RoomManage) Quiet(roomId, uid uint64, expire int64) {
	if roomInfo := r.Get(roomId); roomInfo != nil {
		roomInfo.Quiet.Store(uid, expire)
	}
}

// 是否禁言
func (r *RoomManage) IsQuiet(roomId, uid uint64) bool {
	if roomInfo := r.Get(roomId); roomInfo != nil {
		if expire, ok := roomInfo.Quiet.Load(uid); ok {
			if expire.(int64) > time.Now().Unix() {
				return true
			}
		}
	}

	return false
}

func (r *RoomManage) CancelQuiet(roomId, uid uint64) {
	if roomInfo := r.Get(roomId); roomInfo != nil {
		roomInfo.Quiet.Delete(uid)
	}
}
