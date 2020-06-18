/*
 * @Author: your name
 * @Date: 2020-06-17 16:26:14
 * @LastEditTime: 2020-06-17 16:28:04
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/settings/params.go
 */

package settings

// DecryptSegLen 解密时使用的一次循环下的解密的长度
const DecryptSegLen = 128

// CryptSegLen 加密时使用的一次循环下的解密的长度
const CryptSegLen = 100

// DecryptB64SegLen 解密base64时一次的长度
const DecryptB64SegLen = 172

// LogLev 日志等级
var LogLev = map[string]int{
	"DEBUG": 0,
	"INFO":  1,
	"WARN":  2,
	"ERROR": 3,
	"CRIT":  4,
}
