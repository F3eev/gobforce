package brute

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func (t *Target) MysqlLogin() bool {

	db, err := sql.Open("mysql", t.Username+":"+t.Password+"@tcp("+t.IP+":"+t.Port+")/mysql?charset=utf8")
	if err != nil {
		log.Fatal(err)
		return false
	}
	if err := db.Ping(); err != nil {
		return false
	}
	defer db.Close()
	return true
}

