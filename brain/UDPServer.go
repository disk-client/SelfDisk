/*
 * @Author: xiaoboya
 * @Date: 2020-06-15 09:13:40
 * @LastEditTime: 2020-06-30 11:32:52
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/brain/UDPServer.go
 */

package brain

import (
	"SelfDisk/utils"
	"fmt"
	"net"
	"strconv"
)

// UDPServer UDP打洞的服务端，用以给两个客户端获取IP地址
func UDPServer() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9201})
	if err != nil {
		fmt.Println(err)
		return
	}
	var data = make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error during read: ", err)
			continue
		}
		var content = string(data[:n])
		var clientType = content[:1]
		var tableName string
		switch clientType {
		case "c":
			tableName = "t_client"
		case "s":
			tableName = "t_server"
		default:
			continue
		}
		var username = content[1:]
		var theSQL = `
			select id from t_user where username=$1;
		`
		var uid int
		var theDB = utils.GetConn()
		err = theDB.GetOne(theSQL, []interface{}{username}, []interface{}{&uid})
		if err != nil {
			fmt.Printf("error during auth: no user")
		}
		theSQL = `
			select count(1) from %s where userid=$1;
		`
		theSQL = fmt.Sprintf(theSQL, tableName)
		var isExist int
		theDB.GetOne(theSQL, []interface{}{uid}, []interface{}{&isExist})
		var selfDiskIP = remoteAddr.IP.String()
		var selfDiskPort = strconv.Itoa(remoteAddr.Port)
		if isExist < 1 {
			theSQL = `
				INSERT INTO %s
				(ipaddr, port, userid)
				VALUES($1, $2, $3);
			`
			theSQL = fmt.Sprintf(theSQL, tableName)
			err = theDB.InsertSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
		} else {
			theSQL = `
				UPDATE %s
				SET ipaddr=$1, port=$2
				WHERE userid=$3;
			`
			theSQL = fmt.Sprintf(theSQL, tableName)
			theDB.UpdateSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
		}
		var returnContent = []byte("cmd(addr)" + selfDiskIP + ":9201")
		listener.WriteTo(returnContent, remoteAddr)
		fmt.Println("消息已发送")
	}
}
