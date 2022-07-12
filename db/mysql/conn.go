package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-drive/mysql"
)

var db *sql.DB

func init() {
	sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql,err:" + err.Error())
		os.Exit(1)
	}
}