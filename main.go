package main

import (

	//"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/keithagy-raelyz/GoLive/app"
)

var (
	a *app.App
)

func main() {
	a.StartApp()
}

func resetDBTable(tablename string) {
	db.Exec(fmt.Sprintf("DELETE FROM %s", tablename))
	db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", tablename))
}
