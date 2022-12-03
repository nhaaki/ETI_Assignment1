package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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
	// router.HandleFunc("/api/drive/edit/passenger/{user id}", pEdit).Methods("PUT")
	// router.HandleFunc("/api/drive/edit/driver/{user id}", dEdit).Methods("PUT")
	router.Handle("/api/drive/edit/passenger/{user id}", isAuthorized(pEdit)).Methods("PUT")
	router.Handle("/api/drive/edit/driver/{user id}", isAuthorized(pEdit)).Methods("PUT")
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

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("Invalid Signing Method"))
				}
				aud := "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}
				// verify iss claim
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf(("invalid iss"))
				}

				xmlFile, err := os.Open("../key.xml")
				if err != nil {
					fmt.Println(err)
				}

				defer xmlFile.Close()
				byteValue, _ := ioutil.ReadAll(xmlFile)
				var key Key
				xml.Unmarshal(byteValue, &key)

				var mySigningKey = []byte(key.Value)

				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}
