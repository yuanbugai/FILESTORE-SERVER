package main

import (
	cfg "awesomeProject4/config"
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
	//分块上传接口
	http.HandleFunc("/file/mpupload/init",
		handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",
		handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",
		handler.HTTPInterceptor(handler.CompleteUploadHandler))
	//用户接口
	//http.HandleFunc("/usr/sign/up", handler.Signuphandler)
	http.HandleFunc("/user/signin", handler.Signinhandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserinfoHandler))
	// 初版
	//err := http.ListenAndServe(":80", nil)
	//if err != nil {
	//	fmt.Println("Failed to start server: ", err.Error())
	//}
	fmt.Printf("上传服务启动中，开始监听监听[%s]...\n", cfg.UploadServiceHost)
	// 启动服务并监听端口
	err := http.ListenAndServe(cfg.UploadServiceHost, nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())

	}
}
