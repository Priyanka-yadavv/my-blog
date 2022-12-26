package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	a := App{}
	a.Initialize(DB_USERNAME, DB_PASSWORD, DB_NAME)
	a.Run(":8080")
}
