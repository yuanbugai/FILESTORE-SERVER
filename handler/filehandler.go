package handler

import (
	"awesomeProject4/meta"
	"awesomeProject4/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //GET 方法获取上传主页
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" { //POST 方法获取文件上传内容
		//接受文件流存储到本地
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to get data:", err.Error())
			return
		}

		defer file.Close()

		filePath := "./files/tmp"
		filename := (head.Filename + time.Now().Format("-20060102150405"))
		filename = filepath.Join(filePath, filename)
		err = os.MkdirAll(filePath, os.ModeDir)
		if err != nil {
			fmt.Println("Create directory:", err.Error())
		}
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: filename,
			UploadAt: time.Now().Format("-20060102150405"),
		}
		newFile, err := os.OpenFile(fileMeta.Location, os.O_CREATE|os.O_RDWR, 7777)
		if err != nil {
			fmt.Println(" Failed to create  file", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("Failed to Save file", err.Error())
			return
		}
		//文件哈希值
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetadb(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}

}

//上传已完成
func Uploadsuchandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload successful")
}

//GetFileMetahandle :通过hash获取文件信息
func GetFileMetahandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(filehash)
	fmeta, err2 := meta.GetFileMetadb(filehash)
	if err2 != nil {
		fmt.Println("get file metadata DB error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fmeta)
	if err != nil {
		fmt.Println("get file metadata DB  Write error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
