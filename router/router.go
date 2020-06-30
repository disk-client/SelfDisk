/*
 * @Author: your name
 * @Date: 2020-06-14 22:18:04
 * @LastEditTime: 2020-06-20 17:22:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/router.go
 */

package router

import (
	"SelfDisk/views"

	"github.com/gin-gonic/gin"
)

// MainURL 主要的URL
func MainURL(r *gin.Engine) {
	r.POST("/login", views.Login)
	r.POST("/register", views.Register)

	// downloadURL 下载相关路由
	var downloadURL = r.Group("/download")
	{
		downloadURL.POST("/client", views.DownloadClient)
		downloadURL.POST("/server", views.DownloadServer)
	}

	var clientURL = r.Group("/client")
	{
		clientURL.POST("/server/ipaddr", views.DiskServerAddr)
	}
}
