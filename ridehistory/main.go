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

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Key struct {
	XMLName xml.Name `xml:"key"`
	Value   string   `xml:"value"`
}

func main() {
	router := mux.NewRouter()
	router.Handle("/api/drive/addtohistory", isAuthorized(setHistory)).Methods("POST")
	router.Handle("/api/drive/passenger/gethistory/{user id}", isAuthorized(getHistory)).Methods("GET")
	fmt.Println("Listening at port 6005")
	log.Fatal(http.ListenAndServe(":6005", router))
}

func setHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	var ridePayload map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&ridePayload)
	fmt.Println(ridePayload)

	_, err2 := db.Exec("INSERT INTO RideHistory (driverUID, passengerUID, pcPickUp, pcDropOff) values (?,?,?,?)",
		ridePayload["passengerUID"], ridePayload["driverUID"], ridePayload["pcPickup"], ridePayload["pcDropOff"])
	if err2 != nil {
		w.WriteHeader(http.StatusConflict)
	}

}

func getHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
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
