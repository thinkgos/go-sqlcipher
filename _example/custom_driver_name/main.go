package main

import (
	"database/sql"

	_ "github.com/thinkgos/go-sqlcipher"
)

func main() {
	for _, driver := range sql.Drivers() {
		println(driver)
	}
}
