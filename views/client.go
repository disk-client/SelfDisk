/*
 * @Author: xiaoboya
 * @Date: 2020-06-17 17:52:44
 * @LastEditTime: 2020-06-18 10:16:08
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/views/client.go
 */

package views

import (
	"SelfDisk/utils"

	"github.com/gin-gonic/gin"
)

var execServerTypeDict = map[string]string{
	"windows": "../process/selfdisk_server.exe",
	"linux":   "../process/selfdisk_server_linux",
	"mac":     "../process/selfdisk_server_mac",
}

var execClientTypeDict = map[string]string{
	"windows": "../process/selfdisk_client.exe",
	"linux":   "../process/selfdisk_client_linux",
	"mac":     "../process/selfdisk_client_mac",
}

// DownloadServer 下载服务端
func DownloadServer(c *gin.Context) {
	var descript = utils.RequestDescript{
		Descript: "下载服务端",
		Request:  c,
	}
	var execType = c.Param("exectype")
	v, ok := execServerTypeDict[execType]
	if !ok {
		utils.ReqReturn("目标可执行文件不存在", false, nil, true, descript)
		return
	}
	c.File(v)
}

// DownloadClient 下载客户端
func DownloadClient(c *gin.Context) {
	var descript = utils.RequestDescript{
		Descript: "下载客户端",
		Request:  c,
	}
	var execType = c.Param("exectype")
	v, ok := execClientTypeDict[execType]
	if !ok {
		utils.ReqReturn("目标可执行文件不存在", false, nil, true, descript)
		return
	}
	c.File(v)
}
