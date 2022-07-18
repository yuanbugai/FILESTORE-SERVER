package handler

import (
	dblayer "awesomeProject4/db"
	"awesomeProject4/meta"
	"awesomeProject4/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
		//meta.UpdateFileMeta(fileMeta)
		meta.UpdateFileMetadb(fileMeta)
		r.ParseForm()
		//username := r.Form.Get("username")
		// TODO:根据不同用户名给不同用户添加文件，暂时只能给admin添加文件
		suc := dblayer.OnUserFileUploadFinished("admin", fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		}
		//http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
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

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm, _ := meta.GetFileMetadb(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}
