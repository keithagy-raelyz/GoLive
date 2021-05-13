package main

import (

	//"encoding/json"
	// "fmt"

	"GoLive/app"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v71"

	_ "github.com/go-sql-driver/mysql"
)

var (
	a *app.App
)

func main() {
	godotenv.Load()
	stripe.Key = os.Getenv("STRIPEKEY")
	a = &app.App{}
	a.StartApp()
}

// func resetDBTable(tablename string) {
// 	db.Exec(fmt.Sprintf("DELETE FROM %s", tablename))
// 	db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", tablename))
// }
