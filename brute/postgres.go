package brute

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func (t *Target) PostgresLogin() bool {

	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", t.Username,
		t.Password, t.IP, t.Port, "postgres", "disable")

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		//log.Print("用户名：", username, "    密码: ", password, "      ", "false")
		return false
	}
	if err := db.Ping(); err == nil {
		//log.Print("用户名：", username, "    密码: ", password, "      ", "true")
		defer db.Close()
		return true
	}
	//log.Print("用户名：", username, "    密码: ", password, "      ", "false")
	return false
}
