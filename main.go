package main

import (
	"awesomeProject4/handler"
	"fmt"
	"net/http"
)

func main() {
	//静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//路由协议
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.Uploadsuchandle)
	http.HandleFunc("/file/meta", handler.GetFileMetahandle)
	http.HandleFunc("/usr/sign/up", handler.Signuphandler)
	http.HandleFunc("/user/signin", handler.Signinhandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserinfoHandler))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}

}
