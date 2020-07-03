/*
 * @Author: your name
 * @Date: 2020-07-01 09:15:29
 * @LastEditTime: 2020-07-03 11:20:19
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/brain/relayTCPServer.go
 */

package brain

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

var cache *net.TCPConn = nil

// cacheMap ip:信息
// 但是我现在比较怀疑是不是有不要写个这个变量……
// 可能是我用的地方不太对吧……
// 淦！
var cacheMap = map[string]TCPConnect{}

// makeControl 添加一个tcp端口，用来接收客户端发送的tcp链接
func makeControl() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:8089")
	//打开一个tcp断点监听
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("控制端口已经监听")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		fmt.Println("新的客户端连接到控制端服务进程:" + tcpConn.RemoteAddr().String())
		var content = make([]byte, 1024)
		n, err := tcpConn.Read(content)
		if n == 0 {
			fmt.Println(err)
			continue
		} else {
			var l = strings.Split(tcpConn.RemoteAddr().String(), ":")
			var port, _ = strconv.Atoi(l[1])
			var username = string(content[:n])
			var newtcp = TCPConnect{IP: l[0], Port: port, Name: username, Conn: tcpConn}
			err = newtcp.CheckAuth()
			if err != nil {
				newtcp.Conn.Write(([]byte)("fuck\n"))
			}
			if _, ok := cacheMap[l[0]]; !ok {
				cacheMap[l[0]] = newtcp
				cache = newtcp.Conn
			}
		}
	}
}

//一旦有客户端连接到服务端的话，服务端每隔2秒发送hi消息给到客户端
//如果发送不出去，则认为链路断了，清除cache连接
func control() {
	go func() {
		for {
			for k, v := range cacheMap {
				_, e := v.Conn.Write(([]byte)("hi\n"))
				if e != nil {
					delete(cacheMap, k)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// 接入一个新的链接
func makeAccept() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:8087")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected 8087:" + tcpConn.RemoteAddr().String())
		// 这里是希望去开多个协程，好处理多个转发问题
		go addConnMathAccept(tcpConn, "xiaoboya")
	}
}

// ConnMatch 链接匹配
type ConnMatch struct {
	accept        *net.TCPConn //8087 tcp链路 accept
	acceptAddTime int64        //接受请求的时间
	tunnel        *net.TCPConn //8088 tcp链路 tunnel
}

var connListMap = make(map[string]*ConnMatch)

// 加入匹配（匹配客户端和服务端）
func addConnMathAccept(accept *net.TCPConn, username string) {
	// 这里或者执行到这个函数之前需要进行一次用户认证
	if _, ok := connListMap[username]; !ok {
		connListMap[username] = &ConnMatch{accept, time.Now().Unix(), nil}
		sendMessage("new\n")
	}
}

// 新起链路去链接
func sendMessage(message string) {
	fmt.Println("send Message " + message)
	if cache != nil {
		_, e := cache.Write([]byte(message))
		if e != nil {
			fmt.Println("消息发送异常")
			fmt.Println(e.Error())
		}
	} else {
		fmt.Println("没有客户端连接，无法发送消息")
	}
}

// makeForward 建立第三个tcp接受数据包
func makeForward() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:8088")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected 8088 :" + tcpConn.RemoteAddr().String())
		var l = strings.Split(tcpConn.RemoteAddr().String(), ":")
		var ip = l[0]
		configConnListTunnel(tcpConn, ip)
	}
}

// connListMapUpdate 通道用于更新
var connListMapUpdate = make(chan int)

func configConnListTunnel(tunnel *net.TCPConn, ip string) {
	fmt.Println(ip)
	fmt.Println(cacheMap)
	fmt.Println(connListMap)
	if _, ok := cacheMap[ip]; !ok {
		tunnel.Close()
		fmt.Println("关闭多余的tunnel")
		return
	}
	var name = cacheMap[ip].Name
	v, ok := connListMap[name]
	if !ok || v.tunnel != nil {
		tunnel.Close()
		fmt.Println("关闭多余的tunnel")
		return
	}
	connListMap[name].tunnel = tunnel
	connListMapUpdate <- 1
}

func tcpForward() {
	for {
		select {
		case <-connListMapUpdate:
			for key, connMatch := range connListMap {
				//如果两个都不为空的话，建立隧道连接
				if connMatch.tunnel != nil && connMatch.accept != nil {
					fmt.Println("建立tcpForward隧道连接")
					go joinConn2(connMatch.accept, connMatch.tunnel)
					//从map中删除
					delete(connListMap, key)
				}
			}
		}
	}
}

func joinConn2(conn1 *net.TCPConn, conn2 *net.TCPConn) {
	var f = func(local *net.TCPConn, remote *net.TCPConn) {
		//defer保证close
		// defer local.Close()
		// defer remote.Close()
		//使用io.Copy传输两个tcp连接，
		_, err := io.Copy(local, remote)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("join Conn2 end")
	}
	go f(conn2, conn1)
	go f(conn1, conn2)
}

func releaseConnMatch() {
	for {
		for key, connMatch := range connListMap {
			//如果在指定时间内没有tunnel的话，则释放该连接
			if connMatch.tunnel == nil && connMatch.accept != nil {
				if time.Now().Unix()-connMatch.acceptAddTime > 5 {
					fmt.Println("释放超时连接")
					err := connMatch.accept.Close()
					if err != nil {
						fmt.Println("释放连接的时候出错了:" + err.Error())
					}
					delete(connListMap, key)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

// InitRelayServer 初始化转发服务
func InitRelayServer() {
	// 控制和tcp保持链接
	go control()
	//监听控制端口8009
	go makeControl()
	//监听服务端口8007
	go makeAccept()
	//监听转发端口8008
	go makeForward()
	//定时释放连接
	go releaseConnMatch()
	//执行tcp转发
	go tcpForward()
}
