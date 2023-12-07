package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port1    = 5432
	user     = "test"
	password = "test"
	dbname   = "postgres"
)

type DB struct {
	Id1        string
	Param_id1  string
	Timestamp1 string
	Test       string
}

var (
	tpl  = template.Must(template.ParseFiles("index.html"))
	data = DB{}
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, data)
}

func main() {
	// Подключение
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port1, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.Ping()
	CheckError(err)
	fmt.Println("The database is connected")
	// Запрос
	rows, err := db.Query("SELECT id, param_id, timestamp FROM hour_params LIMIT 50")
	CheckError(err)
	defer rows.Close()
	// Отображение результатов
	for rows.Next() {
		var id, param_id, timestamp string
		err := rows.Scan(&id, &param_id, &timestamp)
		CheckError(err)

		fmt.Printf("id: %s, param_id: %s, timestamp %s\n", id, param_id, timestamp)
		data.Id1 = "<td>" + id + "</td>"
		data.Param_id1 = "<td>" + param_id + "</td>"
		data.Timestamp1 = "<td>" + timestamp + "</td>"
		data.Test = data.Test + "<tr>" + data.Id1 + data.Param_id1 + data.Timestamp1 + "</tr>"

	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// func database() {
// Подключение
// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port1, user, password, dbname)
// db, err := sql.Open("postgres", psqlconn)
// CheckError(err)
// defer db.Close()
// err = db.Ping()
// CheckError(err)
// fmt.Println("The database is connected")
// // Запрос
// rows, err := db.Query("SELECT id, param_id, timestamp FROM hour_params LIMIT 20 WHERE param_id=387")
// CheckError(err)
// defer rows.Close()
// // Отображение результатов
// for rows.Next() {
// 	var id, param_id, timestamp string
// 	err := rows.Scan(&id, &param_id, &timestamp)
// 	CheckError(err)

// 	fmt.Printf("id: %s, param_id: %s, timestamp %s\n", id, param_id, timestamp)
// 	data.Id1 = "<td>" + id + "</td>"
// 	data.Param_id1 = "<td>" + param_id + "</td>"
// 	data.Timestamp1 = "<td>" + timestamp + "</td>"
// 	data.Test = "<tr>" + data.Id1 + data.Param_id1 + data.Timestamp1 + "</tr>" + data.Test

// }
// if err := rows.Err(); err != nil {
// 	panic(err)
// }
// }
