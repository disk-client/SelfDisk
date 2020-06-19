/*
 * @Author: xiaoboya
 * @Date: 2020-06-17 16:15:46
 * @LastEditTime: 2020-06-19 09:29:47
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/utils/crypt.go
 */

package utils

import (
	"SelfDisk/settings"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// PriKey 私钥，与前端加解密相关
var PriKey = new(rsa.PrivateKey)

// PubKey 公钥，与前端加解密相关
var PubKey = new(rsa.PublicKey)

// CryptMsg 加密数据
func CryptMsg(msgIn string) string {
	if PubKey == nil {
		return ""
	}
	var msg = url.QueryEscape(msgIn)
	var enCryptMsg = ""
	var i = 0
	var tLen = len(msg)
	for i < tLen {
		var content []byte
		if tLen > i+settings.CryptSegLen {
			content = []byte(msg[i : i+settings.CryptSegLen])
		} else {
			content = []byte(msg[i:tLen])
		}
		encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, PubKey, content)
		CheckErr(err)
		enCryptMsg += base64.StdEncoding.EncodeToString(encryptBytes)
		i += settings.CryptSegLen
	}
	return enCryptMsg
}

// DecryptMsg 解密数据
func DecryptMsg(msgIn string, key int) string {
	var k = new(rsa.PrivateKey)
	if PriKey == nil {
		return ""
	}
	k = PriKey
	var quoteMsg = ""
	var tLen = len(msgIn)
	var i = 0
	msgIn = strings.Replace(msgIn, " ", "", -1)
	for i < tLen {
		enc, err := base64.StdEncoding.DecodeString(msgIn[i : i+settings.DecryptB64SegLen])
		CheckErr(err)
		res, err := rsa.DecryptPKCS1v15(rand.Reader, k, enc)
		CheckErr(err)
		quoteMsg += string(res)
		i += settings.DecryptB64SegLen
	}
	result, err := url.QueryUnescape(quoteMsg)
	CheckErr(err)
	return result
}

// ReqReturn 请求返回，加密response
func ReqReturn(msg string, succ bool, data interface{}, noCrypt bool, req RequestDescript) {
	var succSign = 1

	if succ == false {
		succSign = 0
	}
	if req.Request != nil {
		var lev = settings.LogLev["ERROR"]
		if succ == true {
			lev = settings.LogLev["INFO"]
		}
		var infomsg = ""
		if msg != "" {
			infomsg = ":" + msg
		}
		fmt.Println(lev, infomsg)
		// var action = fmt.Sprintf("成功:%d", succSign)
		// UserLog(req, action+infomsg)
	}
	var res = map[string]interface{}{
		"succ": succSign,
		"info": msg,
		"data": data,
	}
	jsonRes, err := json.Marshal(res)
	CheckErr(err)
	if settings.DeBug {
		req.Request.JSON(200, gin.H{
			"data": res,
		})
	} else {
		if noCrypt == true {
			req.Request.JSON(200, gin.H{
				"data": "0" + string(jsonRes),
			})
		} else {
			var endRes = CryptMsg(string(jsonRes))
			req.Request.JSON(200, gin.H{
				"data": "1" + endRes,
			})
		}
	}
}

// ReqParse 接收请求，解密request
func ReqParse(c *gin.Context, obj interface{}) {
	// body, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	c.String(400, err.Error())
	// 	c.Abort()
	// }
	// c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	var msg = c.PostForm("msg")
	// var msg = string(body)
	msg = DecryptMsg(msg, 1)
	var err = json.Unmarshal([]byte(msg), obj)
	CheckErr(err)
}

// InitKeys 初始化公钥和私钥
func InitKeys() {
	var f []byte
	var err error

	// 生成解密用的私钥
	if CheckFileIsExist("../cert/fontPrivate.pem") {
		f, err = ioutil.ReadFile("../cert/fontPrivate.pem")
		CheckErr(err)
	} else {
		return
	}
	block, _ := pem.Decode(f)
	if block == nil {
		return
	}
	PriKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	if PriKey == nil {
		return
	}

	// 生成加密用的公钥
	if CheckFileIsExist("../cert/backendPublic.pem") {
		f, err = ioutil.ReadFile("../cert/backendPublic.pem")
		CheckErr(err)
	} else {
		return
	}
	block, _ = pem.Decode(f)
	if block == nil {
		return
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if pubKey == nil {
		CheckErr(err)
		return
	}
	value, ok := pubKey.(*rsa.PublicKey)
	if ok == true {
		PubKey = value
	} else {
		return
	}
}
