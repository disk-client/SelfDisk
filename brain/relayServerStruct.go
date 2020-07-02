/*
 * @Author: your name
 * @Date: 2020-07-01 11:11:05
 * @LastEditTime: 2020-07-02 09:33:46
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/brain/relayServerStruct.go
 */

package brain

import (
	"SelfDisk/utils"
	"errors"
	"fmt"
	"net"
	"strconv"
)

// TCPConnect 定义一个TCP链接结构用于存储链接信息
type TCPConnect struct {
	IP   string
	Port int
	Name string
	Conn *net.TCPConn
}

// GetAddr 获取地址
func (tc *TCPConnect) GetAddr() string {
	var portStr = strconv.Itoa(tc.Port)
	return fmt.Sprintf("%s:%s", tc.IP, portStr)
}

// CheckAuth 核验身份
func (tc *TCPConnect) CheckAuth() error {
	var theSQL = `
		select count(1) from t_user where username=$1;
	`
	var n int
	var theDB = utils.GetConn()
	var err = theDB.GetOne(theSQL, []interface{}{tc.Name}, []interface{}{&n})
	if err != nil || n == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
