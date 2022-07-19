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

	//文件接口
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.Uploadsuchandle))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetahandle))
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	//用户接口
	http.HandleFunc("/usr/sign/up", handler.Signuphandler)
	http.HandleFunc("/user/signin", handler.Signinhandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserinfoHandler))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}

}
