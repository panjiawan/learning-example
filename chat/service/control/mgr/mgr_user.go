package mgr

import (
	"errors"
	"github.com/coldwind/artist/pkg/ihttp"
	"github.com/coldwind/artist/pkg/ilog"
	"github.com/panjiawan/note/chat/service/control/common"
	"go.uber.org/zap"
	"sync"
)

type UserManage struct {
	sync.RWMutex
	userPool   map[uint64]map[int64]*ihttp.WSClient
	noAuthPool map[int64]*ihttp.WSClient
}

var UserMgr *UserManage

func init() {
	UserMgr = &UserManage{
		userPool:   make(map[uint64]map[int64]*ihttp.WSClient),
		noAuthPool: make(map[int64]*ihttp.WSClient),
	}
}

func (m *UserManage) Add(client *ihttp.WSClient) error {
	userData := client.UserData.(*common.UserData)

	m.Lock()
	defer m.Unlock()

	if !userData.IsAuth() {
		m.noAuthPool[client.ConnId] = client
		ilog.Debug("add to guest", zap.Int64("connid", client.ConnId))
	} else {
		if _, ok := m.userPool[userData.Uid]; !ok {
			m.userPool[userData.Uid] = make(map[int64]*ihttp.WSClient)
		}

		m.userPool[userData.Uid][client.ConnId] = client
		ilog.Debug("add to auth", zap.Int64("connid", client.ConnId), zap.Uint64("uid", userData.Uid))
	}
	return nil
}

func (m *UserManage) Del(client *ihttp.WSClient) error {
	userData := client.UserData.(*common.UserData)

	m.Lock()
	defer m.Unlock()

	if !userData.IsAuth() {
		delete(m.noAuthPool, client.ConnId)
		ilog.Debug("remove conn", zap.Int64("connid", client.ConnId))
	} else if _, ok := m.userPool[userData.Uid]; ok {
		delete(m.userPool[userData.Uid], client.ConnId)

		if len(m.userPool[userData.Uid]) == 0 {
			delete(m.userPool, userData.Uid)
			ilog.Debug("remove authed uid", zap.Uint64("uid", userData.Uid))
		}

		ilog.Debug("remove authed conn", zap.Uint64("uid", userData.Uid), zap.Int64("connid", client.ConnId))

	}

	return nil
}

// 删除用户UID下的所有连接
func (m *UserManage) DelByUid(uid uint64) error {
	m.Lock()
	defer m.Unlock()
	delete(m.userPool, uid)

	return nil
}

// 根据connid删除 只能删除未认证连接
func (m *UserManage) DelByConnId(connId int64) error {
	m.Lock()
	defer m.Unlock()
	delete(m.noAuthPool, connId)

	return nil
}

// 未认证的用户 完成认证
func (m *UserManage) Authed(client *ihttp.WSClient) error {
	userData := client.UserData.(*common.UserData)

	m.Lock()
	defer m.Unlock()

	if !userData.IsAuth() {
		m.noAuthPool[client.ConnId] = client
		return errors.New("auth failure")
	} else {
		delete(m.noAuthPool, client.ConnId)
		ilog.Debug("del from no auth pool", zap.Int64("connid", client.ConnId))

		if _, ok := m.userPool[userData.Uid]; !ok {
			m.userPool[userData.Uid] = make(map[int64]*ihttp.WSClient)
		}

		m.userPool[userData.Uid][client.ConnId] = client
		ilog.Debug("add to auth pool", zap.Int64("connid", client.ConnId))

	}
	return nil
}

// 发送给指定用户
func (m *UserManage) SendToUser(uid uint64, msg []byte) error {
	m.RLock()
	defer m.RUnlock()
	if user, ok := m.userPool[uid]; ok {
		for _, v := range user {
			v.Send(msg)
		}
	}

	return nil
}

// 广播到消息
func (m *UserManage) Broadcast(msg []byte) error {
	m.RLock()
	defer m.RUnlock()
	for _, user := range m.userPool {
		for _, conn := range user {
			conn.Send(msg)
		}
	}

	for _, conn := range m.noAuthPool {
		conn.Send(msg)
	}

	return nil
}

func (m *UserManage) BroadcastWithoutRoom(roomId uint64, msg []byte) error {
	m.RLock()
	defer m.RUnlock()
	for _, user := range m.userPool {
		for _, conn := range user {
			userdata := conn.UserData.(*common.UserData)
			if userdata.RoomId == roomId {
				continue
			}

			conn.Send(msg)
		}
	}

	for _, conn := range m.noAuthPool {
		conn.Send(msg)
	}

	return nil
}

func (m *UserManage) GetRoomIdsByUid(uid uint64) []uint64 {
	roomIds := make([]uint64, 0, 2)
	m.RLock()
	defer m.RUnlock()
	if user, ok := m.userPool[uid]; ok {
		for _, v := range user {
			userdata := v.UserData.(*common.UserData)
			if userdata.RoomId > 0 {
				roomIds = append(roomIds, userdata.RoomId)
			}
		}
	}

	return roomIds
}
