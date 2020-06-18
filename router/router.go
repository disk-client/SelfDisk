/*
 * @Author: your name
 * @Date: 2020-06-14 22:18:04
 * @LastEditTime: 2020-06-18 10:46:59
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
}
