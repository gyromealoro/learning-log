package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Data struct {
	ID   int    `json:"id"`
	Note string `json:"note"`
}

var db *sql.DB

func Open_data_base() {
	godotenv.Load()
	var err error

	Db_host := os.Getenv("DB_HOST")
	Db_port := os.Getenv("DB_PORT")
	Db_user := os.Getenv("DB_USER")
	Db_pass := os.Getenv("DB_PASSWORD")
	Db_name := os.Getenv("DB_NAME")

	connect := "host=" + Db_host + " port=" + Db_port + " user=" + Db_user + " password=" + Db_pass + " dbname=" + Db_name

	db, err = sql.Open("postgres", connect)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to the database!")

}

func Addtask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	s_name := r.FormValue("s_name")
	if s_name == "" {
		http.Error(w, "Note cannot be empty", http.StatusBadRequest)
		return
	}
	//Decode the json data to struct and store it into data variable

	//Insert the data to the database
	_, err := db.Exec("insert into test_t(name) values($1)", s_name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusCreated)
	//fmt.Fprintf(w, "task is added to the database!")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Deletedata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	data := r.FormValue("data")

	_, err := db.Exec("delete from test_t where id = $1", data)

	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func ReadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return

	}
	rows, err := db.Query("select id, name FROM test_t")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var notes []Data
	for rows.Next() {
		var d Data
		rows.Scan(&d.ID, &d.Note)
		notes = append(notes, d)

	}

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, notes)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(notes)

}

func main() {

	Open_data_base()
	defer db.Close()

	http.HandleFunc("/del", Deletedata)
	http.HandleFunc("/add", Addtask)
	http.HandleFunc("/", ReadData)
	http.ListenAndServe(":3000", nil)

}
