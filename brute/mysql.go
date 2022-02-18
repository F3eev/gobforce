package brute

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (t *Target) MysqlLogin() bool {
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)

	db, err := sql.Open("mysql", t.Username+":"+t.Password+"@tcp("+t.IP+":"+t.Port+")/mysql?charset=utf8")
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	if err := db.Ping(); err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	defer db.Close()
	return true
}
