package db

import (
	mydb "awesomeProject4/db/mysql"
	"fmt"
)

type token struct {
	Token string
}

func GenTokenbyusername(username string) *token {
	token := token{}
	stmt, err := mydb.Dbconnect().Prepare(
		"select user_token from tbl_user_token where user_name=? limit 1",
	)
	if err != nil {
		fmt.Println("error hanpend in token.go +", err.Error())
	}
	row := stmt.QueryRow(username)
	row.Scan(&token.Token)
	fmt.Println(token.Token)
	return &token
}
