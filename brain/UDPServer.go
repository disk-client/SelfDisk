/*
 * @Author: your name
 * @Date: 2020-06-15 09:13:40
 * @LastEditTime: 2020-06-18 10:43:58
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
			fmt.Printf("error during read: %s", err)
		}
		var username = string(data[:n])
		var theSQL = `
			select id from t_user where username=$1;
		`
		var uid int
		err = diskutils.TheDB.GetOne(theSQL, []interface{}{username}, []interface{}{&uid})
		if err != nil {

		}
		theSQL = `
			select count(1) from t_server where userid=$1;
		`
		var isExist int
		diskutils.TheDB.GetOne(theSQL, []interface{}{uid}, []interface{}{&isExist})
		var selfDiskIP = remoteAddr.IP
		var selfDiskPort = remoteAddr.Port
		if isExist > 1 {
			theSQL = `
				INSERT INTO public.t_server
				(ipaddr, port, userid)
				VALUES($1, $2, $3);			
			`
			utils.TheDB.InsertSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
		} else {
			theSQL = `
				UPDATE public.t_server
				SET ipaddr=$1, port=$2
				WHERE userid=$3;
			`
			utils.TheDB.UpdateSQL(theSQL, []interface{}{selfDiskIP, selfDiskPort, uid})
		}
	}
}
