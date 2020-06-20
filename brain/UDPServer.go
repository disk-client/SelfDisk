/*
 * @Author: xiaoboya
 * @Date: 2020-06-15 09:13:40
 * @LastEditTime: 2020-06-20 09:33:17
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/brain/UDPServer.go
 */

package brain

import (
	"SelfDisk/utils"
	diskutils "SelfDisk/utils"
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
		}
		var content = string(data[:n])
		fmt.Println(content)
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
		err = diskutils.TheDB.GetOne(theSQL, []interface{}{username}, []interface{}{&uid})
		if err != nil {
			fmt.Printf("error during auth: no user")
		}
		theSQL = `
			select count(1) from %s where userid=$1;
		`
		theSQL = fmt.Sprintf(theSQL, tableName)
		var isExist int
		diskutils.TheDB.GetOne(theSQL, []interface{}{uid}, []interface{}{&isExist})
		var selfDiskIP = remoteAddr.IP.String()
		var selfDiskPort = strconv.Itoa(remoteAddr.Port)
		if isExist < 1 {
			theSQL = `
				INSERT INTO %s
				(ipaddr, port, userid)
				VALUES($1, $2, $3);
			`
			theSQL = fmt.Sprintf(theSQL, tableName)
			fmt.Println(selfDiskIP, selfDiskPort)
			err = utils.TheDB.InsertSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
			fmt.Println(err)
		} else {
			theSQL = `
				UPDATE %s
				SET ipaddr=$1, port=$2
				WHERE userid=$3;
			`
			theSQL = fmt.Sprintf(theSQL, tableName)
			utils.TheDB.UpdateSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
		}
	}
}
