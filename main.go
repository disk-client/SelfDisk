/*
 * @Author: your name
 * @Date: 2020-06-13 21:49:06
 * @LastEditTime: 2020-06-18 11:21:50
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/main.go
 */

package main

import (
	"SelfDisk/brain"
	"SelfDisk/router"
	"SelfDisk/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	go brain.UDPServer()
	// 生成公私钥
	utils.InitKeys()
	var r = gin.Default()
	router.MainURL(r)
	r.Run(":8080")
}
