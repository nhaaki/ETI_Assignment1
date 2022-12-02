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

type Driver struct {
	UserID       int    `json:"User ID"`
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	FirstName    string `json:"First Name"`
	LastName     string `json:"Last Name"`
	MobileNo     string `json:"Mobile Number"`
	EmailAddress string `json:"Email Address"`
	IdNo         string `json:"Identification Number"`
	CarLicenseNo string `json:"Car License Number"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/ride/driver/{user id}", drv)
	router.HandleFunc("/api/ride/driver", dlogin)
	fmt.Println("Listening at port 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}

func drv(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)
	if r.Method == "POST" {
		w.Header().Set("Content-type", "application/json")
		d := json.NewDecoder(r.Body)
		var t Driver
		t.UserID = 999 //placeholder value for account creation
		err := d.Decode(&t)
		if err != nil {
			panic(err)
		} else {
			var count int
			db.QueryRow("Select count(*) from Drivers where Username=?", t.Username).Scan(&count)
			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, "Username taken. Try another username.")
			} else {
				w.WriteHeader(http.StatusAccepted)
				_, err := db.Exec("INSERT INTO Drivers (UserID, Username, Password, FirstName, LastName, MobileNo, EmailAddress, IdNo, CarLicenseNo) values (?,?,?,?,?,?,?,?,?)",
					t.UserID, t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress, t.IdNo, t.CarLicenseNo)
				if err != nil {
					panic(err.Error())
				}
				fmt.Fprintf(w, "Account with username %s created!", t.Username)
			}
		}
	} else if r.Method == "PUT" {
		w.Header().Set("Content-type", "application/json")
		d := json.NewDecoder(r.Body)
		var t Driver
		x, _ := strconv.Atoi(params["user id"])
		t.UserID = x
		err := d.Decode(&t)
		if err != nil {
			panic(err)
		} else {
			var count int
			db.QueryRow("Select count(*) from Drivers where Username=? and UserID!=?", t.Username, t.UserID).Scan(&count)
			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprint(w, "Username taken. Try another username.")
			} else {
				w.WriteHeader(http.StatusAccepted)
				_, err := db.Exec("UPDATE Drivers SET Username=?, Password=?, FirstName=?, LastName=?, MobileNo=?, EmailAddress=?, IdNo=?, CarLicenseNo=? WHERE UserID=?",
					t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress, t.IdNo, t.CarLicenseNo, t.UserID)
				if err != nil {
					panic(err.Error())
				}
				fmt.Fprintf(w, "Account with username %s updated!", t.Username)
			}
		}
	}

}

func dlogin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	Username := r.URL.Query().Get("username")
	Password := r.URL.Query().Get("password")

	var count int
	db.QueryRow("Select count(*) from Drivers where Username=? and Password=?", Username, Password).Scan(&count)
	if count != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No user found or wrong password.")
	} else {
		w.WriteHeader(http.StatusAccepted)
		var d Driver
		db.QueryRow("select * from Drivers where Username=? and Password=?", Username, Password).Scan(&d.UserID,
			&d.Username, &d.Password, &d.FirstName, &d.LastName, &d.MobileNo, &d.EmailAddress, &d.IdNo, &d.CarLicenseNo)

		res, _ := json.MarshalIndent(d, "", "\t")
		fmt.Fprintf(w, string(res))
	}

}
