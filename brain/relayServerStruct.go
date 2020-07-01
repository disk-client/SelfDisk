/*
 * @Author: your name
 * @Date: 2020-07-01 11:11:05
 * @LastEditTime: 2020-07-01 11:11:20
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/brain/relayServerStruct.go
 */

package brain

import (
	"fmt"
	"strconv"
)

// TCPConnect 定义一个TCP链接结构用于存储链接信息
type TCPConnect struct {
	IP   string
	Port int
	Name string
}

// GetAddr 获取地址
func (tc *TCPConnect) GetAddr() string {
	var portStr = strconv.Itoa(tc.Port)
	return fmt.Sprintf("%s:%s", tc.IP, portStr)
}

// CheckAuth 核验身份
func (tc *TCPConnect) CheckAuth() {

}
