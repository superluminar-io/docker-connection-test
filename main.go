package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"database/sql"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	DbHost string `env:"DB_HOST"`
	DbUser string `env:"DB_USERNAME"`
	DbPass string `env:"DB_PASSWORD"`
}

func main() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		if err != nil {
			return
		}
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", cfg.DbUser, cfg.DbPass, cfg.DbHost))
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		// Prepare statement for reading data
		stmtOut, err := db.Prepare(`SELECT 1`)
		if err != nil {
			panic(err.Error())
		}
		defer stmtOut.Close()
		_, err = stmtOut.Exec()
		if err != nil {
			panic(err.Error())
		}

		_, err = fmt.Fprintf(w, "DB Connection successful!")
		if err != nil {
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
