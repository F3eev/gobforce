package brute

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func (t *Target) PostgresLogin() bool {

	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", t.Username,
		t.Password, t.IP, t.Port, "postgres", "disable")
	debugMessage := fmt.Sprintf("%s %s %s:%s", RunFileName(), RunFuncName(), t.IP, t.Port)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		t.Logrus.Debug(fmt.Sprintf("%s %v", debugMessage, err))
		return false
	}
	if err := db.Ping(); err == nil {

		defer db.Close()
		return true
	}
	return false
}
