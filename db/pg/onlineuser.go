package pg

import (
	"encoding/json"
	"fmt"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

type OnlineUser struct {
	User       Person `json:"user"`
	TokenID    string `json:"id"`
	ClientType string `json:"clienttype"`
	FromIp     string `json:"fromip"`
	ExpireTime int64  `json:"expiretime"`
}

// OnlineUser.Add 增加在线用户
func (ou *OnlineUser) Add() (resStatus i18n.ResKey, err error) {
	key := fmt.Sprintf("%s%s%d", ou.ClientType, ":", ou.User.ID)
	jsonL, _ := json.Marshal(ou)
	err = cache.SetOther(key, jsonL)
	return
}

// OnlineUser.Get 获取在线用户
func (ou *OnlineUser) Get() (exist int32, resStatus i18n.ResKey, err error) {
	key := fmt.Sprintf("%s%s%d", ou.ClientType, ":", ou.User.ID)
	exist, v, err := cache.GetOther(key)
	if exist == 1 {
		err = json.Unmarshal(v, &ou)
	}
	return
}

// OnlineUser.Del 删除在线用户
func (ou *OnlineUser) Del() (resStatus i18n.ResKey, err error) {
	key := fmt.Sprintf("%s%s%d", ou.ClientType, ":", ou.User.ID)
	err = cache.DelOther(key)
	if err != nil {
		resStatus = i18n.CodeInternalError
		return
	}
	return
}

// GetAllOnlineUser 获取所有在线用户列表
func GetAllOnlineUser() (ous []OnlineUser, resStatus i18n.ResKey, err error) {
	ous = make([]OnlineUser, 0)
	//获取用户表中的所有用户
	sqlStr := `select id from sysuser where dr=0`

	rows, err := db.Query(sqlStr)
	if err != nil {
		resStatus = i18n.CodeInternalError
		zap.L().Error("GetAllOnlineUser db.Query failed", zap.Error(err))
		return
	}
	defer rows.Close()

	//提取数据
	for rows.Next() {
		var ou OnlineUser
		err = rows.Scan(&ou.User.ID)
		if err != nil {
			resStatus = i18n.CodeInternalError
			zap.L().Error("GetAllOnlineUser row.Next() failed", zap.Error(err))
			return
		}

		//从缓存中获取当前登录用户信息
		for _, clientType := range pub.ValidClientTypes {
			ou.ClientType = clientType
			exist, resStatus, errGet := ou.Get()
			if errGet != nil {
				zap.L().Error("GetAllOnlineUser ou.Get failed:", zap.Error(errGet))
				return ous, resStatus, errGet
			}
			if exist == 1 && ou.ExpireTime > time.Now().Unix() {
				ous = append(ous, ou)
			}
		}
	}
	return
}
