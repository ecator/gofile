package server

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ecator/gofile/util"
)

// 从cookie中获取token并重新设置过期时间，没有的话就重新设置一个
func getToken(w http.ResponseWriter, r *http.Request) string {
	name := "token"
	cookie, err := r.Cookie(name)
	token := ""
	if err == nil {
		token = cookie.Value
	} else {
		token = util.Md5FromStr(strconv.Itoa(rand.Intn(1000)) + time.Now().String() + r.UserAgent())
	}
	cookie = new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Path = "/"
	cookie.Expires = time.Now().AddDate(1, 1, 1)
	http.SetCookie(w, cookie)
	return token
}
