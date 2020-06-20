/*
 * @Author: your name
 * @Date: 2020-06-20 16:15:59
 * @LastEditTime: 2020-06-20 16:38:12
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/args/aboutClient.go
 */

package args

import (
	"SelfDisk/utils"
	"errors"
)

// ClientCmd 客户端指令
type ClientCmd struct {
	Name string
	Cmd  string
}

// CheckUser 检查用户
func (info *ClientCmd) CheckUser() int {
	var theDB = utils.GetConn()
	var theSQL = `
		select id from t_user where username=$1;
	`
	var userID int
	var err = theDB.GetOne(theSQL, []interface{}{info.Name}, []interface{}{&userID})
	if err != nil {
		return 0
	}
	return userID
}

// ExecCmd 执行客户端指令
func (info *ClientCmd) ExecCmd(uID int) (res interface{}) {
	switch info.Cmd {
	case "ServerAddr":
		res = info.getServerAddr(uID)
	}
	return
}

func (info *ClientCmd) getServerAddr(uID int) interface{} {
	var theDB = utils.GetConn()
	var theSQL = `
		select ipaddr,port from t_server where userid=$1;
	`
	var ip, port string
	var err = theDB.GetOne(theSQL, []interface{}{uID}, []interface{}{&ip, &port})
	if err != nil {
		return errors.New("客户端不存在")
	}
	return []string{ip, port}
}
