package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Key struct {
	XMLName xml.Name `xml:"key"`
	Value   string   `xml:"value"`
}

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
	router.HandleFunc("/api/drive/login/passenger", plogin)
	router.HandleFunc("/api/drive/login/driver", dlogin)
	fmt.Println("Listening at port 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}

// Login function for PASSENGER users
func plogin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var userValues map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&userValues)

	var count int
	db.QueryRow("Select count(*) from Passengers where Username=? and Password=?", userValues["Username"], userValues["Password"]).Scan(&count)
	if count != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No user found or wrong password.")
	} else {
		w.WriteHeader(http.StatusAccepted)
		var p Passenger
		db.QueryRow("select * from Passengers where Username=? and Password=?", userValues["Username"], userValues["Password"]).Scan(&p.UserID,
			&p.Username, &p.Password, &p.FirstName, &p.LastName, &p.MobileNo, &p.EmailAddress)

		res, _ := json.MarshalIndent(p, "", "\t")
		fmt.Fprintf(w, string(res))
	}
}

// Login function for DRIVER users
func dlogin(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var userValues map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&userValues)

	var count int
	db.QueryRow("Select count(*) from Drivers where Username=? and Password=?", userValues["Username"], userValues["Password"]).Scan(&count)
	if count != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No user found or wrong password.")
	} else {
		w.WriteHeader(http.StatusAccepted)
		var d Driver
		db.QueryRow("select * from Drivers where Username=? and Password=?", userValues["Username"], userValues["Password"]).Scan(&d.UserID,
			&d.Username, &d.Password, &d.FirstName, &d.LastName, &d.MobileNo, &d.EmailAddress, &d.IdNo, &d.CarLicenseNo)
		db.Exec("INSERT INTO LiveRides (driverUID, status) values (?,?)",
			d.UserID, "Available")

		res, _ := json.MarshalIndent(d, "", "\t")
		fmt.Fprintf(w, string(res))
	}

}
