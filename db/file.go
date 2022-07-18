package db

import (
	mydb "awesomeProject4/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := mydb.Dbconnect().Prepare(
		"insert ignore into tab_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type Tablefile struct {
	Filehash string `json:"filehash"`
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
	Fileaddr string `json:"fileaddr"`
}

// GetFileMeta:从mysql获取文件源信息
func GetFileMeta(filehash string) (*Tablefile, error) {

	stmt, err := mydb.Dbconnect().Prepare("select file_sha1,file_addr,file_name,file_size from tab_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
	}
	tfile := Tablefile{}
	stmt.QueryRow(filehash).Scan(&tfile.Filehash, &tfile.Fileaddr, &tfile.Filename, &tfile.Filesize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}
