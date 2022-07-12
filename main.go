package main

import (
	mysql "awesomeProject4/db/mysql"
	"awesomeProject4/handler"
	"fmt"
	"net/http"
)

func main() {
	//路由协议
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.Uploadsuchandle)
	http.HandleFunc("file/meta", handler.GetFileMetahandle)
	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}
	dbconnect := mysql.Dbconnect()
	if dbconnect != nil {
		fmt.Println("successful")
	}

}
