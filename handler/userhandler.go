package handler

import (
	dblayer "awesomeProject4/db"
	"awesomeProject4/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salted = "*#890"
)

//Signuphandler:处理用户注册
func Signuphandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/test.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("404")
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid paramter"))
		fmt.Println("invalid paramter")
		return
	}
	enc_passwd := util.Sha1([]byte(passwd + pwd_salted))
	suc := dblayer.Usersingup(username, enc_passwd)
	if suc {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("Filed"))
	}

}

//Signinhandler :登录接口
func Signinhandler(w http.ResponseWriter, r *http.Request) {
	//校验用户名及密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salted))
	pwdchecked := dblayer.Usersignin(username, encPassword)
	if !pwdchecked {
		w.Write([]byte("failed"))
		return
	}
	//生成Token访问
	token := GenToken(username)
	upres := dblayer.Updatetoken(username, token)
	if !upres {
		w.Write([]byte("failed"))
		return
	}
	//TODO:登录成功后重定向到首页

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

//UserinfoHandler :查询用户信息
func UserinfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	////2.验证token是否有效
	//valid := IsTokenValid(token)
	//if !valid {
	//	fmt.Println("token invalid")
	//
	//}

	//3.查询用户信息
	user, err := dblayer.GetuserinfoByuer_name(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())

}

//////////////////////////////
/////////////////////////////
/////////////////////////////

func GenToken(username string) string {
	// (username+times+Token_sqlt)+times[:8]
	ts := fmt.Sprintf("%s", time.Now())
	tokenprefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenprefix + ts[:8]
}
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
