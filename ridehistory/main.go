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
	router.Handle("/api/drive/addtohistory", isAuthorized(setHistory)).Methods("POST")
	router.Handle("/api/drive/passenger/gethistory/{user id}", isAuthorized(getHistory)).Methods("GET")
	fmt.Println("Listening at port 6005")
	log.Fatal(http.ListenAndServe(":6005", router))
}

// This function adds a record to the RideHistory table. This function is called after a ride has ended.
func setHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveDataDB")
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
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

// This function gets all the record under a single passenger's ID. This function is called through their menu, and returns
// the trips in reverse chronological order in a map.
func getHistory(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveDataDB")
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
		var pcPickUp string
		var pcDropOff string
		var rideDate string
		var d Driver
		err = results.Scan(&passengerUID, &driverUID, &pcPickUp, &pcDropOff, &rideDate)
		if err != nil {
			panic(err.Error())
		}

		client := &http.Client{}
		url := "http://localhost:5000/api/drive/get/driver?userid=" + strconv.Itoa(driverUID)
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			req.Header.Set("Token", r.Header["Token"][0])
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 404 || res.StatusCode == 409 {
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(body))
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					err := json.Unmarshal(body, &d)
					if err != nil {
						panic(err.Error())
					}
				}
			}
		}

		rideObj := map[string]string{
			"driverFirstName": d.FirstName,
			"driverLastName":  d.LastName,
			"carLicenseNo":    d.CarLicenseNo,
			"pcPickup":        pcPickUp,
			"pcDropOff":       pcDropOff,
			"rideDate":        rideDate,
		}
		allRides[x] = rideObj

		x++
	}

	if len(allRides) < 1 {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No ride history found...")
	} else {
		w.Header().Set("Content-type", "application/json")
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
