/*
 * @Author: xiaoboya
 * @Date: 2020-06-15 09:13:40
 * @LastEditTime: 2020-07-01 09:12:34
 * @LastEditors: Please set LastEditors
 * @Description: 早先尝试的UDP式的点对点直连是一个不错的想法，但是最终发现效果查强人意
                现在已经尝试了基于tcp的转发。虽然受服务器带宽限制比较大，但是我觉得还是可以接受的
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
		sendMsg(selfDiskIP, "cmd(addr)"+selfDiskIP+":9201", remoteAddr.Port)
		fmt.Println("消息已发送")
	}
}

func sendMsg(host, msg string, port int) {
	var srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: 9202} // 注意端口必须固定
	dstAddr := &net.UDPAddr{IP: net.ParseIP(host), Port: port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	conn.Write([]byte(msg))
	conn.Close()
}
