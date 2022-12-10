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
	router.HandleFunc("/api/drive/create/passenger", createPassengerUser).Methods("POST")
	router.HandleFunc("/api/drive/create/driver", createDriverUser).Methods("POST")
	router.HandleFunc("/api/drive/checkusername", checkUsername).Methods("POST")
	router.HandleFunc("/api/drive/login/passenger", plogin)
	router.HandleFunc("/api/drive/login/driver", dlogin)
	router.Handle("/api/drive/get/passenger", isAuthorized(getPDetails)).Methods("GET")
	router.Handle("/api/drive/get/driver", isAuthorized(getDDetails)).Methods("GET")
	router.Handle("/api/drive/edit/passenger/{user id}", isAuthorized(pEdit)).Methods("PUT")
	router.Handle("/api/drive/edit/driver/{user id}", isAuthorized(dEdit)).Methods("PUT")
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
			_, err := db.Exec("INSERT INTO Passengers (Username, Password, FirstName, LastName, MobileNo, EmailAddress) values (?,?,?,?,?,?)",
				t.Username, t.Password, t.FirstName, t.LastName, t.MobileNo, t.EmailAddress)
			if err != nil {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "Error - Information entered may be in the wrong format.")
			} else {
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, "Account with username %s created!", t.Username)
			}
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
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "Error - Information entered may be in the wrong format.")
			} else {
				fmt.Fprintf(w, "Account with username %s created!", t.Username)
			}
		}
	}

}

func checkUsername(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No username input.")
	} else {
		var count int
		username := r.URL.Query().Get("username")
		accType := r.URL.Query().Get("type")

		if accType == "Passenger" {
			db.QueryRow("Select count(*) from Passengers where Username=?", username).Scan(&count)
			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "Error - Username taken. Try another username.")

			} else {
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, "Username is free.")
			}
		} else {
			db.QueryRow("Select count(*) from Drivers where Username=?", username).Scan(&count)
			if count != 0 {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintf(w, "Error - Username taken. Try another username.")

			} else {
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, "Username is free.")
			}
		}

	}

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

		client := &http.Client{}
		url := "http://localhost:6002/api/drive/driver/setstatus?userid=" + strconv.Itoa(d.UserID)
		if req, err := http.NewRequest("POST", url, nil); err == nil {
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 202 {
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(body))
				}
			}
		}

		res, _ := json.MarshalIndent(d, "", "\t")
		fmt.Fprintf(w, string(res))
	}
}

func getPDetails(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No user ID input.")
	} else {
		userID := r.URL.Query().Get("userid")
		var p Passenger

		err = db.QueryRow("Select * from Passengers where UserID=?", userID).Scan(&p.UserID, &p.Username, &p.Password, &p.FirstName, &p.LastName, &p.MobileNo, &p.EmailAddress)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Error - Unable to retrieve information.")
			panic(err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			res, _ := json.MarshalIndent(p, "", "\t")
			fmt.Fprintf(w, string(res))
		}
	}
}

func getDDetails(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	querystringmap := r.URL.Query()
	if len(querystringmap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error - No user ID input.")
	} else {
		userID := r.URL.Query().Get("userid")
		var d Driver

		err = db.QueryRow("Select * from Drivers where UserID=?", userID).Scan(&d.UserID, &d.Username, &d.Password, &d.FirstName, &d.LastName, &d.MobileNo, &d.EmailAddress, &d.IdNo, &d.CarLicenseNo)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Error - Unable to retrieve information.")
			panic(err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			res, _ := json.MarshalIndent(d, "", "\t")
			fmt.Fprintf(w, string(res))
		}
	}
}

func pEdit(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	params := mux.Vars(r)

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
			var prevUN string
			var prevPW string
			db.QueryRow("Select Username, Password from Passengers where UserID=?", t.UserID).Scan(&prevUN, &prevPW)
			unpwPayload := make(map[string]string)
			unpwPayload["Username"] = prevUN
			unpwPayload["Password"] = prevPW
			res, _ := json.MarshalIndent(unpwPayload, "", "\t")
			fmt.Fprintf(w, string(res))
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
			var prevUN string
			var prevPW string
			db.QueryRow("Select Username, Password from Drivers where UserID=?", t.UserID).Scan(&prevUN, &prevPW)
			unpwPayload := make(map[string]string)
			unpwPayload["Username"] = prevUN
			unpwPayload["Password"] = prevPW
			res, _ := json.MarshalIndent(unpwPayload, "", "\t")
			fmt.Fprintf(w, string(res))
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
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("You're unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}

		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}
