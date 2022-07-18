package db

import (
	mydb "awesomeProject4/db/mysql"
	"fmt"
)

//Usersingup: 通过用户名和密码判断是否能注册这个用户
func Usersingup(username string, password string) bool {
	stmt, err := mydb.Dbconnect().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`)values (?,?)")
	if err != nil {
		fmt.Println("Failed to prepare insert,err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert,err:" + err.Error())

	}
	affected, err := ret.RowsAffected()
	if err != nil && affected > 0 {
		return true
	}
	return false
}

//Usersignin : 通过用户名判断是否有这个用户
func Usersignin(username, enc_passwd string) bool {
	stmt, err := mydb.Dbconnect().Prepare("select * from tbl_user where user_name =? limit 1")
	if err != nil {
		fmt.Println("can't search " + err.Error())
		return false
	}
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("can't query " + err.Error())
	} else if rows == nil {
		fmt.Println("username not found")
		return false
	}
	//查询数据转换map类型数组
	prow := mydb.ParseRows(rows)
	if len(prow) > 0 && string(prow[0]["user_pwd"].([]byte)) == enc_passwd {
		return true
	}
	return false
}

//Updatetoken 刷新用户的token
func Updatetoken(username string, token string) bool {
	stmt, err := mydb.Dbconnect().Prepare("replace into tbl_user_token(`user_name`,`user_token`)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true

}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	lastActiveAt string
	Status       int
}

func GetuserinfoByuer_name(username string) (User, error) {
	user := User{}
	stmt, err := mydb.Dbconnect().Prepare(
		"select user_name,signup_at from tbl_user where  user_name=? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
