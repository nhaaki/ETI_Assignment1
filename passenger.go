package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Passenger struct {
	UserID       int    `json:"User ID"`
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	FirstName    string `json:"First Name"`
	LastName     string `json:"Last Name"`
	MobileNo     string `json:"Mobile Number"`
	EmailAddress string `json:"Email Address"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/ride/passenger/{user id}", psg)
	router.HandleFunc("/api/ride/passenger", plogin)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func psg(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	params := mux.Vars(r)
	if r.Method == "POST" {
		w.Header().Set("Content-type", "application/json")
		d := json.NewDecoder(r.Body)
		var t Passenger
		t.UserID = 999 //placeholder value for account creation
		err := d.Decode(&t)
		if err != nil {
			panic(err)
		} else {
			var count int
			db.QueryRow("Select count(*) from Passengers where Username=?", t.Username).Scan(&count)

			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, "Username taken. Try another username.")
			} else {
				w.WriteHeader(http.StatusAccepted)
				_, err := db.Exec("INSERT INTO Passengers (Username, Password, FirstName, LastName, MobileNo, EmailAddress) values (?,?,?,?,?,?)",
					t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress)
				if err != nil {
					panic(err.Error())
				}
				fmt.Fprintf(w, "Account with username %s created!", t.Username)
			}
		}
	} else if r.Method == "PUT" {
		w.Header().Set("Content-type", "application/json")
		d := json.NewDecoder(r.Body)
		var t Passenger
		x, _ := strconv.Atoi(params["user id"])
		t.UserID = x
		err := d.Decode(&t)
		if err != nil {
			panic(err)
		} else {
			var count int
			db.QueryRow("Select count(*) from Passengers where Username=? and UserID!=?", t.Username, t.UserID).Scan(&count)
			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, "Username taken. Try another username.")
			} else {
				w.WriteHeader(http.StatusAccepted)
				_, err := db.Exec("UPDATE Passengers SET Username=?, Password=?, FirstName=?, LastName=?, MobileNo=?, EmailAddress=? WHERE UserID=?",
					t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress, t.UserID)
				if err != nil {
					panic(err.Error())
				}
				fmt.Fprintf(w, "Account with username %s updated!", t.Username)
			}
		}
	}

}

func plogin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	Username := r.URL.Query().Get("username")
	Password := r.URL.Query().Get("password")

	var count int
	db.QueryRow("Select count(*) from Passengers where Username=? and Password=?", Username, Password).Scan(&count)
	if count != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No user found or wrong password.")
	} else {
		w.WriteHeader(http.StatusAccepted)
		var p Passenger
		db.QueryRow("select * from Passengers where Username=? and Password=?", Username, Password).Scan(&p.UserID,
			&p.Username, &p.Password, &p.FirstName, &p.LastName, &p.MobileNo, &p.EmailAddress)

		res, _ := json.MarshalIndent(p, "", "\t")
		fmt.Fprintf(w, string(res))
	}

}
