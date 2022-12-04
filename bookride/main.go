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

func main() {
	router := mux.NewRouter()
	router.Handle("/api/drive/passenger/book/{user id}", isAuthorized(bookRide)).Methods("POST")
	fmt.Println("Listening at port 6002")
	log.Fatal(http.ListenAndServe(":6002", router))
}

func bookRide(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user id"])

	var pcValues map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&pcValues)

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	db.QueryRow("Select count(*) from LiveRides where status=?", "Available").Scan(&count)

	if count != 0 {
		_, err := db.Exec("UPDATE LiveRides SET passengerUID=?, pcPickUp=?, pcDropOff=?, status=? where status=? limit 1", userID,
			pcValues["pcPickUp"], pcValues["pcDropOff"], "Assigned")
		if err != nil {
			panic(err.Error())
		}
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "Rider found!")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Error - No riders available...")
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
