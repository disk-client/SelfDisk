/*
 * @Author: xiaoboya
 * @Date: 2020-06-13 21:49:06
 * @LastEditTime: 2020-07-01 10:42:55
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/main.go
 */

package main

import (
	"SelfDisk/brain"
	"SelfDisk/router"
	"SelfDisk/utils"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	go brain.InitRelayServer()
	// 生成公私钥
	utils.InitKeys()
	var r = gin.Default()

	// 设置session
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(10 * time.Minute), //10min
		Path:   "/",
	})
	r.Use(sessions.Sessions("mysession", store))

	router.MainURL(r)
	r.Run(":8080")
}
