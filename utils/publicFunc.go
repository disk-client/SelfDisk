/*
 * @Author: 肖博雅
 * @Date: 2020-05-09 17:53:10
 * @LastEditTime: 2020-06-17 16:51:47
 * @LastEditors: Please set LastEditors
 * @Description: 没有特殊归类的通用函数
 * @FilePath: /ZeroTrust/utils/publicFunc.go
 */

package utils

import (
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestDescript 请求描述，描述+gin.Context
type RequestDescript struct {
	Descript string
	Request  *gin.Context
}

// CheckErr 判断错误，有错误就是直接结束这个协程
func CheckErr(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

// InsertSlice 对切片进行插入操作
func InsertSlice(slice, insertion []interface{}, index int) []interface{} {
	return append(slice[:index], append(insertion, slice[index:]...)...)
}

// HandleSQL 用户处理拼接出来的sql语句（由于pg比较特殊，需要使用$1-$n这样的数字进行占位）
func HandleSQL(sql string, params []interface{}) string {
	for i := range params {
		var n = strconv.Itoa(i + 1)
		sql = strings.Replace(sql, "$n", "$"+n, 1)
	}
	return sql
}

// B2S 把[]uint8类型转换成字符串，也是用于处理Postgres返回类型为元素为字符串的数组的部分
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, b)
	}
	var a = string(ba)
	a = strings.Replace(a, "{", "", -1)
	a = strings.Replace(a, "}", "", -1)
	return a
}

// GetTimes 处理审计部分的时间的
func GetTimes(datetype string) (today, tomorrow, startTime, datetypeCN string) {
	var nowTime = time.Now()
	today = nowTime.Format("2006-01-02")
	tomorrow = nowTime.AddDate(0, 0, 1).Format("2006-01-02")
	if datetype == "week" {
		startTime = nowTime.AddDate(0, 0, -7).Format("2006-01-02")
		datetypeCN = "近七天"
	} else if datetype == "month" {
		startTime = nowTime.AddDate(0, -1, 0).Format("2006-01-02")
		datetypeCN = "近一月"
	} else if datetype == "quarter" {
		startTime = nowTime.AddDate(0, -3, 0).Format("2006-01-02")
		datetypeCN = "近三月"
	} else if datetype == "halfyear" {
		startTime = nowTime.AddDate(0, -6, 0).Format("2006-01-02")
		datetypeCN = "近半年"
	} else if datetype == "year" {
		startTime = nowTime.AddDate(-1, 0, 0).Format("2006-01-02")
		datetypeCN = "近一年"
	} else {
		startTime = nowTime.AddDate(0, 0, -7).Format("2006-01-02")
		datetypeCN = "近七天"
	}
	return today, tomorrow, startTime, datetypeCN
}

// CheckFileIsExist 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// Rounding 四舍五入
func Rounding(f float64) int {
	f = f + 0.5
	var res = int(math.Floor(f))
	return res
}

// VerifyEmailFormat 验证电子邮箱数据格式
func VerifyEmailFormat(email string) bool {
	var pattern = `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	var reg = regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyMobileFormat 验证电子邮箱数据格式
func VerifyMobileFormat(phone string) bool {
	var regular = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	var reg = regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

// VerifyUsernameOrPasswordFormat 验证用户名和密码格式
func VerifyUsernameOrPasswordFormat(username string) bool {
	var ruler = "[a-z0-9A-Z]+$"
	var reg = regexp.MustCompile(ruler)
	return reg.MatchString(username)
}
