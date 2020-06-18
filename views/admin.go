/*
 * @Author: your name
 * @Date: 2020-06-15 21:12:03
 * @LastEditTime: 2020-06-18 16:03:31
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/views/admin.go
 */

package views

import (
	"SelfDisk/args"
	"SelfDisk/settings"
	"SelfDisk/utils"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var descript = utils.RequestDescript{
		Descript: "注册新用户",
		Request:  c,
	}
	var info args.RegisterParams
	if settings.DeBug {
		c.BindJSON(&info)
	} else {
		utils.ReqParse(c, &info)
	}
	if msg := info.Check(); msg != "" {
		utils.ReqReturn(msg, false, nil, false, descript)
		return
	}
	var err = info.SaveToDB()
	if err != nil {
		utils.ReqReturn("新建数据失败", false, nil, false, descript)
		return
	}
	var session = sessions.Default(c)
	session.Set("username", info.Name)
	session.Save()
	utils.ReqReturn("注册用户成功", true, nil, false, descript)
	return
}

// Login 用户登录
func Login(c *gin.Context) {
	var descript = utils.RequestDescript{
		Descript: "用户登录",
		Request:  c,
	}
	var info args.LoginParams
	if settings.DeBug {
		c.BindJSON(&info)
	} else {
		utils.ReqParse(c, &info)
	}
	if err := info.Check(); err != nil {
		utils.ReqReturn("登录失败", false, nil, false, descript)
		return
	}
	var session = sessions.Default(c)
	session.Set("username", info.Name)
	session.Save()
	utils.ReqReturn("登录成功", true, nil, false, descript)
	return
}
