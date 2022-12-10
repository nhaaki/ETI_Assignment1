package main

import (
	"context"
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
	router.HandleFunc("/api/drive/checkactiveride", checkActiveRide).Methods("GET")
	router.HandleFunc("/api/drive/getdriverstatus", getDriverStatus).Methods("GET")
	router.Handle("/api/drive/getridedetails", isAuthorized(getRideDetails)).Methods("GET")
	router.HandleFunc("/api/drive/driver/setstatus", setStatus).Methods("POST")
	router.Handle("/api/drive/driver/setstatus", isAuthorized(setStatus)).Methods("DELETE")
	router.Handle("/api/drive/passenger/assigndriver/{user id}", isAuthorized(assignDriver)).Methods("POST")
	router.Handle("/api/drive/driver/ridefunctions/{user id}", isAuthorized(rideFunctions)).Methods("POST")
	fmt.Println("Listening at port 6002")
	log.Fatal(http.ListenAndServe(":6002", router))
}

func checkActiveRide(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No UserID input.")
	} else {
		userID := r.URL.Query().Get("userid")
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var count int
		db.QueryRow("Select count(*) from LiveRides where passengerUID=?", userID).Scan(&count)
		fmt.Fprint(w, count)
	}
}

func getDriverStatus(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No UserID input.")
	} else {
		driverUID := r.URL.Query().Get("userid")
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var status string
		db.QueryRow("Select status from LiveRides where driverUID=?", driverUID).Scan(&status)
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, status)
	}
}

func getRideDetails(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No UserID input.")
	} else {
		driverUID := r.URL.Query().Get("userid")
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()
		var passengerUID string
		var pcPickUp string
		var pcDropOff string
		var p Passenger

		db.QueryRow("Select l.passengerUID, l.pcPickUp, l.pcDropOff from LiveRides l where l.driverUID=?", driverUID).Scan(&passengerUID, &pcPickUp, &pcDropOff)

		client := &http.Client{}
		url := "http://localhost:5000/api/drive/get/passenger?userid=" + passengerUID
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			req.Header.Set("Token", r.Header["Token"][0])
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 404 || res.StatusCode == 409 {
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(body))
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					err := json.Unmarshal(body, &p)
					if err != nil {
						panic(err.Error())
					}
				}
			}
		}

		ridePayload := map[string]string{
			"pFirstName": p.FirstName,
			"pLastName":  p.LastName,
			"pUID":       passengerUID,
			"pcPickUp":   pcPickUp,
			"pcDropOff":  pcDropOff,
		}
		res, _ := json.MarshalIndent(ridePayload, "", "\t")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, string(res))
	}
}

func setStatus(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No UserID input.")
	} else {
		driverUID := r.URL.Query().Get("userid")
		db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		if r.Method == "POST" {
			db.Exec("INSERT INTO LiveRides (driverUID, status) values (?,?)",
				driverUID, "Available")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "Record added successfully.")
		} else if r.Method == "DELETE" {
			db.Exec("DELETE FROM LiveRides where driverUID=?", driverUID)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, "Successfully logged out!")
		}
	}
}

func assignDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user id"])

	var pcValues map[string]string
	d := json.NewDecoder(r.Body)
	d.Decode(&pcValues)

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var count int
	db.QueryRow("Select count(*) from LiveRides where status=?", "Available").Scan(&count)

	if count != 0 {
		_, err := db.Exec("UPDATE LiveRides SET passengerUID=?, pcPickUp=?, pcDropOff=?, status=? where status=? limit 1", userID,
			pcValues["pcPickup"], pcValues["pcDropoff"], "Assigned", "Available")
		if err != nil {
			panic(err.Error())
		}
		var driverUID string
		db.QueryRow("Select driverUID from LiveRides where passengerUID=?", userID).Scan(&driverUID)
		var d Driver

		client := &http.Client{}
		url := "http://localhost:5000/api/drive/get/driver?userid=" + driverUID
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

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "=======================\nRider found!")
		fmt.Fprintln(w, "Name: "+d.FirstName+" "+d.LastName)
		fmt.Fprintln(w, "Car license number: "+d.CarLicenseNo)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Error - No riders available...")
	}

}

func rideFunctions(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveFunctionDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	params := mux.Vars(r)
	pUID, _ := strconv.Atoi(params["user id"])

	ctx, cancel := context.WithCancel(context.Background())

	router := mux.NewRouter()
	router.Handle("/api/drive/startride/{user id}", isAuthorized(func(h http.ResponseWriter, z *http.Request) {

		params := mux.Vars(z)
		dUID, _ := strconv.Atoi(params["user id"])

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "=======================\nRide started!")
		fmt.Fprintln(w, "   -           __")
		fmt.Fprintln(w, " --          ~( @\\   \\")
		fmt.Fprintln(w, "---   _________]_[__/_>________")
		fmt.Fprintln(w, "     /  ____ \\ <>     |  ____  \\")
		fmt.Fprintln(w, "    =\\_/ __ \\_\\_______|_/ __ \\__D")
		fmt.Fprintln(w, "________(__)_____________(__)____")

		_, err := db.Exec("UPDATE LiveRides SET status=? where passengerUID=? and driverUID=?", "Ongoing", pUID, dUID)
		if err != nil {
			panic(err.Error())
		}
		cancel()
	})).Methods("POST")

	router.Handle("/api/drive/endride/{user id}", isAuthorized(func(h http.ResponseWriter, z *http.Request) {
		params := mux.Vars(z)
		dUID, _ := strconv.Atoi(params["user id"])

		var pcPickUp string
		var pcDropOff string
		db.QueryRow("Select pcPickUp, pcDropOff from LiveRides where passengerUID=?", pUID).Scan(&pcPickUp, &pcDropOff)

		db.Exec("UPDATE LiveRides SET passengerUID=?,pcPickUp=?,pcDropOff=?, status=? where passengerUID=? and driverUID=?",
			nil, nil, nil, "Available", pUID, dUID)
		histPayload := map[string]string{
			"passengerUID": strconv.Itoa(pUID),
			"driverUID":    strconv.Itoa(dUID),
			"pcPickup":     pcPickUp,
			"pcDropOff":    pcDropOff,
		}
		res, _ := json.MarshalIndent(histPayload, "", "\t")
		h.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(h, string(res))
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "Ride ended!\n=======================")

		cancel()
	})).Methods("POST")

	router.Handle("/api/drive/cancelride/{user id}", isAuthorized(func(h http.ResponseWriter, z *http.Request) {
		params := mux.Vars(z)
		dUID, _ := strconv.Atoi(params["user id"])
		db.Exec("UPDATE LiveRides SET passengerUID=?,pcPickUp=?,pcDropOff=?, status=? where driverUID=?",
			nil, nil, nil, "Available", dUID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Driver cancelled the ride.\n=======================")
		cancel()
	})).Methods("POST")

	srv := &http.Server{
		Addr:    "0.0.0.0:6003",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	err2 := srv.Shutdown(context.Background())
	if err2 != nil {
		log.Println(err2)
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
