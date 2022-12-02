package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	router.HandleFunc("/api/drive/create/passenger", createPassengerUser).Methods("POST")
	router.HandleFunc("/api/drive/create/driver", createDriverUser)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func createPassengerUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	d := json.NewDecoder(r.Body)
	var t Passenger

	err2 := d.Decode(&t)
	if err2 != nil {
		panic(err2)
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

}

func createDriverUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	w.Header().Set("Content-type", "application/json")
	d := json.NewDecoder(r.Body)
	var t Driver
	err2 := d.Decode(&t)
	if err2 != nil {
		panic(err2)
	} else {
		var count int
		db.QueryRow("Select count(*) from Drivers where Username=?", t.Username).Scan(&count)
		if count != 0 {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "Username taken. Try another username.")
		} else {
			w.WriteHeader(http.StatusAccepted)
			_, err := db.Exec("INSERT INTO Drivers (Username, Password, FirstName, LastName, MobileNo, EmailAddress, IdNo, CarLicenseNo) values (?,?,?,?,?,?,?,?)",
				t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress, t.IdNo, t.CarLicenseNo)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, "Account with username %s created!", t.Username)
		}
	}

}
