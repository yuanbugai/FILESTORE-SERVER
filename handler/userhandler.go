package handler

import (
	dblayer "awesomeProject4/db"
	"awesomeProject4/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	pwd_salted = "*#890"
)

//Signuphandler:处理用户注册
func Signuphandler(c *gin.Context) {

	c.Redirect(http.StatusOK, "/static/view/signup.html")
}

//Dosignuphandler:处理post请求
func Dosignuphandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")
	if len(username) < 3 || len(passwd) < 5 {

		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid paramter",
			"code": -1,
		})
		fmt.Println("invalid paramter")
		return
	}
	encPassword := util.Sha1([]byte(passwd + pwd_salted))
	pwdchecked := dblayer.Usersignin(username, encPassword)
	if !pwdchecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid signin",
			"code": -2,
		})
		return
	} else {
		fmt.Println("login successful")
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success login",
			"code": 0,
		})
	}
}

//Signinhandler :登录接口
func Signinhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	//校验用户名及密码
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salted))
	pwdchecked := dblayer.Usersignin(username, encPassword)
	if !pwdchecked {
		w.Write([]byte("failed"))
		return
	} else {
		fmt.Println("login successful")
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
func Dosigninhandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	encPassword := util.Sha1([]byte(password + pwd_salted))
	pwdchecked := dblayer.Usersignin(username, encPassword)
	if !pwdchecked {
		c.JSON(http.StatusOK, gin.H{
			"Msg":  "signup failed",
			"code": -1,
		})
		return
	}
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
