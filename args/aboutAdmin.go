/*
 * @Author: 肖博雅
 * @Date: 2020-06-15 21:14:59
 * @LastEditTime: 2020-06-20 16:21:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/args/aboutAdmin.go
 */

package args

import (
	"SelfDisk/utils"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

// RegisterParams 注册参数
type RegisterParams struct {
	Email      string
	Phone      string
	Name       string
	Password   string
	RePassword string
}

// Check 检查是否符合注册要求
func (info *RegisterParams) Check() (msg string) {
	if info.Password != info.RePassword {
		return "两次密码不一致"
	}
	// 用户名只包含大小写和数字
	if ok := utils.VerifyUsernameOrPasswordFormat(info.Name); !ok {
		return "用户名只包含大小写和数字"
	}
	// 密码只包含大小写和数字
	if ok := utils.VerifyUsernameOrPasswordFormat(info.Password); !ok {
		return "用户名只包含大小写和数字"
	}
	// 手机号正则验证
	if ok := utils.VerifyMobileFormat(info.Phone); !ok {
		return "手机号不合规"
	}
	// 邮箱正则验证
	if ok := utils.VerifyEmailFormat(info.Email); !ok {
		return "邮箱不合规"
	}
	return ""
}

// SaveToDB 保存至数据库
func (info *RegisterParams) SaveToDB() error {
	var theSQL = `
		INSERT INTO public.t_user
		(username, "password", email, phone)
		VALUES($1, $2, $3, $4);	
	`
	var cred = "self" + info.Password + "disk"
	var hash = md5.New()
	var credByte = []byte(cred)
	hash.Write(credByte)
	cred = hex.EncodeToString(hash.Sum(nil))
	var theDB = utils.GetConn()
	theDB.InsertSQL(theSQL, []interface{}{info.Name, cred, info.Email, info.Phone})
	return nil
}

// LoginParams 登录参数
type LoginParams struct {
	Name     string
	Password string
}

// Check 验证登录信息
func (info *LoginParams) Check() error {
	var cred = "self" + info.Password + "disk"
	var hash = md5.New()
	var credByte = []byte(cred)
	hash.Write(credByte)
	cred = hex.EncodeToString(hash.Sum(nil))
	fmt.Println(cred)
	var theSQL = `
		select count(1) from t_user where username=$1 and password=$2;
	`
	var n int
	var theDB = utils.GetConn()
	var err = theDB.GetOne(theSQL, []interface{}{info.Name, cred}, []interface{}{&n})
	if err != nil || n == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
