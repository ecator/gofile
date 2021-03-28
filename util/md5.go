package util

import (
	"crypto/md5"
	"encoding/hex"
)

// 计算字符串的MD5值
func Md5FromStr(str string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(str))
	return hex.EncodeToString(algorithm.Sum(nil))
}
