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
	router.HandleFunc("/api/drive/edit/passenger/{user id}", pEdit).Methods("PUT")
	router.HandleFunc("/api/drive/edit/driver/{user id}", dEdit).Methods("PUT")
	fmt.Println("Listening at port 6001")
	log.Fatal(http.ListenAndServe(":6001", router))
}

func pEdit(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)
	w.Header().Set("Content-type", "application/json")

	var t Passenger
	x, _ := strconv.Atoi(params["user id"])
	t.UserID = x

	d := json.NewDecoder(r.Body)
	err2 := d.Decode(&t)
	if err2 != nil {
		panic(err2)
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

func dEdit(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	d := json.NewDecoder(r.Body)
	params := mux.Vars(r)

	var t Driver
	x, _ := strconv.Atoi(params["user id"])
	t.UserID = x

	err2 := d.Decode(&t)
	if err2 != nil {
		panic(err2)
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
