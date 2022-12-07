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

func main() {
	router := mux.NewRouter()
	router.Handle("/api/drive/addtohistory", isAuthorized(setHistory)).Methods("POST")
	router.Handle("/api/drive/passenger/gethistory/{user id}", isAuthorized(getHistory)).Methods("GET")
	fmt.Println("Listening at port 6005")
	log.Fatal(http.ListenAndServe(":6005", router))
}

func setHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ridePayload map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&ridePayload)

	_, err2 := db.Exec("INSERT INTO RideHistory (driverUID, passengerUID, pcPickUp, pcDropOff) values (?,?,?,?)",
		ridePayload["driverUID"], ridePayload["passengerUID"], ridePayload["pcPickup"], ridePayload["pcDropOff"])
	if err2 != nil {
		w.WriteHeader(http.StatusConflict)
	}

}

func getHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user id"])

	results, err2 := db.Query("select * from RideHistory where passengerUID=? ORDER BY rideDate DESC", userID)
	if err2 != nil {
		panic(err2.Error())
	}
	allRides := make(map[int]map[string]string)
	x := 1
	for results.Next() {
		var passengerUID string
		var driverUID int
		var dFirstName string
		var dLastName string
		var CarLicenseNo string
		var pcPickUp string
		var pcDropOff string
		var rideDate string
		err = results.Scan(&passengerUID, &driverUID, &pcPickUp, &pcDropOff, &rideDate)
		if err != nil {
			panic(err.Error())
		}

		db.QueryRow("Select FirstName, LastName, CarLicenseNo from Drivers where UserID=?", driverUID).Scan(&dFirstName, &dLastName, &CarLicenseNo)

		rideObj := map[string]string{
			"driverFirstName": dFirstName,
			"driverLastName":  dLastName,
			"carLicenseNo":    CarLicenseNo,
			"pcPickup":        pcPickUp,
			"pcDropOff":       pcDropOff,
			"rideDate":        rideDate,
		}
		allRides[x] = rideObj

		x++
	}

	if len(allRides) < 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No ride history found...")
	} else {
		res, _ := json.MarshalIndent(allRides, "", "\t")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(res))
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
