package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

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
	var currentPassenger Passenger
	var currentDriver Driver

	for {
		var choice string

		fmt.Println("\n=======================")
		fmt.Println("|| Welcome to DRIVE! ||")
		fmt.Println("=======================")
		fmt.Println("1. Sign in")
		fmt.Println("2. Create account")
		fmt.Println("3. Exit")
		fmt.Print("Enter an option:")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
		signInLoop:
			for {

				var Username string
				var Password string
				fmt.Println("\n=================")
				fmt.Println("1. Sign in as a Passenger")
				fmt.Println("2. Sign in as a Driver")
				fmt.Println("3. Go back")
				fmt.Print("Enter an option:")
				fmt.Scanln(&choice)

				switch choice {
				case "1":
					// Sign in as a passenger
					fmt.Println("\n=================")
					fmt.Print("Enter your username: ")
					fmt.Scanln(&Username)
					fmt.Print("Enter your password: ")
					fmt.Scanln(&Password)

					client := &http.Client{}
					url := "http://localhost:8082/get/JWT"
					if req, err := http.NewRequest("GET", url, nil); err == nil {
						if res, err := client.Do(req); err == nil {
							defer res.Body.Close()
							if res.StatusCode == 409 {
								fmt.Printf("Error in retrieving JWT token... \n")
							} else if res.StatusCode == 202 {

								body, _ := ioutil.ReadAll(res.Body)
								t := string(body)

								url = "http://localhost:6000/api/drive/login/passenger?username=" + Username + "&password=" + Password
								if req, err := http.NewRequest("GET", url, nil); err == nil {
									req.Header.Set("Token", t)
									if res, err := client.Do(req); err == nil {
										defer res.Body.Close()
										if res.StatusCode == 404 {
											fmt.Printf("Error - user not found or incorrect password! \n")
										} else if res.StatusCode == 202 {
											fmt.Printf("Logging in... \n\n")
											body, err := ioutil.ReadAll(res.Body)
											var p Passenger
											if err != nil {
												panic(err)
											} else {
												err := json.Unmarshal(body, &p)
												if err != nil {
													panic(err)
												}
												currentPassenger = p
												PassengerFunctions(&currentPassenger)
												break signInLoop
											}
										}
									}
								}
							}
						}
					}

				case "2":
					// Sign in as a driver
					fmt.Println("\n=================")
					fmt.Print("Enter your username: ")
					fmt.Scanln(&Username)
					fmt.Print("Enter your password: ")
					fmt.Scanln(&Password)

					client := &http.Client{}
					url := "http://localhost:8082/get/JWT"
					if req, err := http.NewRequest("GET", url, nil); err == nil {
						if res, err := client.Do(req); err == nil {
							defer res.Body.Close()
							if res.StatusCode == 409 {
								fmt.Printf("Error in retrieving JWT token... \n")
							} else if res.StatusCode == 202 {

								body, _ := ioutil.ReadAll(res.Body)
								t := string(body)

								url = "http://localhost:6000/api/drive/login/driver?username=" + Username + "&password=" + Password
								if req, err := http.NewRequest("GET", url, nil); err == nil {
									req.Header.Set("Token", t)
									if res, err := client.Do(req); err == nil {
										defer res.Body.Close()
										if res.StatusCode == 404 {
											fmt.Printf("Error - user not found or incorrect password! \n\n")
										} else if res.StatusCode == 202 {
											fmt.Printf("Logging in... \n\n")
											body, err := ioutil.ReadAll(res.Body)
											var d Driver
											if err != nil {
												panic(err)
											} else {
												err := json.Unmarshal(body, &d)
												if err != nil {
													panic(err)
												}
												currentDriver = d
												DriverFunctions(&currentDriver)
												break signInLoop
											}
										}
									}
								}
							}
						}
					}
				case "3":
					break signInLoop
				default:
					fmt.Println("\n=================")
					fmt.Printf("Invalid option! Try again.")
					continue
				}
			}

		case "2":
			var Username string
			var Password string
			var FirstName string
			var LastName string
			var MobileNo string
			var EmailAddress string

		createAccountLoop:
			for {
				fmt.Println("\n=================")
				fmt.Println("1. Create Passenger account")
				fmt.Println("2. Create Driver account")
				fmt.Println("3. Go back")
				fmt.Print("Enter an option:")
				fmt.Scanln(&choice)

				switch choice {
				case "1":

				psgloop:
					for {
						fmt.Println("\n=================")
						fmt.Print("Enter your username: ")
						fmt.Scanln(&Username)

						var count int
						db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
						if err != nil {
							panic(err.Error())
						}
						defer db.Close()

						db.QueryRow("Select count(*) from Passengers where Username=?", Username).Scan(&count)
						if count != 0 {
							fmt.Println("Username taken. Try another username.")
							continue
						}
						fmt.Print("Enter your password: ")
						fmt.Scanln(&Password)
						fmt.Print("Enter your first name: ")
						fmt.Scanln(&FirstName)
						fmt.Print("Enter your last name: ")
						fmt.Scanln(&LastName)
						fmt.Print("Enter your mobile no: ")
						fmt.Scanln(&MobileNo)
						fmt.Print("Enter your email address: ")
						fmt.Scanln(&EmailAddress)

						newAccount := Passenger{999, Username, Password, FirstName, LastName, MobileNo, EmailAddress}
						url := "http://localhost:5000/api/drive/create/passenger"
						postBody, _ := json.Marshal(newAccount)
						resBody := bytes.NewBuffer(postBody)
						client := &http.Client{}
						if req, err := http.NewRequest("POST", url, resBody); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 202 {
									fmt.Printf("Account %s created successfully!\n", newAccount.Username)
									break psgloop

								} else {
									fmt.Println("Error in account creation...")
									break psgloop
								}
							}
						}
					}
					break createAccountLoop
				case "2":
					var IdNo string
					var CarLicenseNo string
				drvloop:
					for {
						fmt.Println("\n=================")
						fmt.Print("Enter your username: ")
						fmt.Scanln(&Username)

						var count int
						db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
						if err != nil {
							panic(err.Error())
						}
						defer db.Close()

						db.QueryRow("Select count(*) from Drivers where Username=?", Username).Scan(&count)
						if count != 0 {
							fmt.Println("Username taken. Try another username.")
							continue
						}
						fmt.Print("Enter your password: ")
						fmt.Scanln(&Password)
						fmt.Print("Enter your first name: ")
						fmt.Scanln(&FirstName)
						fmt.Print("Enter your last name: ")
						fmt.Scanln(&LastName)
						fmt.Print("Enter your mobile no: ")
						fmt.Scanln(&MobileNo)
						fmt.Print("Enter your email address: ")
						fmt.Scanln(&EmailAddress)
						fmt.Print("Enter your ID number: ")
						fmt.Scanln(&IdNo)
						fmt.Print("Enter your car license number: ")
						fmt.Scanln(&CarLicenseNo)

						newAccount := Driver{999, Username, Password, FirstName, LastName, MobileNo, EmailAddress, IdNo, CarLicenseNo}
						url := "http://localhost:5000/api/drive/create/driver"
						postBody, _ := json.Marshal(newAccount)
						resBody := bytes.NewBuffer(postBody)
						client := &http.Client{}
						if req, err := http.NewRequest("POST", url, resBody); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 202 {
									fmt.Printf("Account %s created successfully!\n", newAccount.Username)
									break drvloop
								} else {
									fmt.Printf("Error in account creation...")
									break drvloop
								}
							}
						}
					}
					break createAccountLoop
				case "3":
					break createAccountLoop
				default:
					fmt.Println("\n=================")
					fmt.Printf("Invalid option! Try again.")
					continue
				}
			}

		case "3":
			fmt.Println("\n=================")
			fmt.Println("See you later!")
			os.Exit(3)
		default:
			fmt.Println("\n=================")
			fmt.Println("Invalid option! Try again.")
		}
	}
}

func PassengerFunctions(curp *Passenger) {
	var choice string

psgloop:
	for {
		fmt.Println("\n=======================")
		fmt.Println("|| Passenger Functions ||")
		fmt.Println("=======================")
		fmt.Printf("Hello, %s %s! \n", curp.FirstName, curp.LastName)
		fmt.Println("=======================\n")
		fmt.Println("1. Book a trip")
		fmt.Println("2. Retrieve all trips")
		fmt.Println("3. Edit account information")
		fmt.Println("4. Exit to main menu")
		fmt.Print("Enter an option:")
		fmt.Scanln(&choice)

		switch choice {
		case "3":
			var Username string
			var Password string
			var FirstName string
			var LastName string
			var MobileNo string
			var EmailAddress string
		editDetailsLoop:
			for {
				fmt.Println("\n=================")
				fmt.Println("1. Edit username and password")
				fmt.Println("2. Edit first name and last name")
				fmt.Println("3. Edit mobile number")
				fmt.Println("4. Edit email address")
				fmt.Println("5. Exit edit menu")
				fmt.Print("Enter an option:")
				fmt.Scanln(&choice)

				switch choice {
				case "1":
					fmt.Print("Enter new username: ")
					fmt.Scanln(&Username)
					fmt.Print("Enter new password: ")
					fmt.Scanln(&Password)

					curp.Username = Username
					curp.Password = Password

					editPsg(curp)
				case "2":
					fmt.Print("Enter new first name: ")
					fmt.Scanln(&FirstName)
					fmt.Print("Enter new last name: ")
					fmt.Scanln(&LastName)

					curp.FirstName = FirstName
					curp.LastName = LastName

					editPsg(curp)
				case "3":
					fmt.Print("Enter new mobile number: ")
					fmt.Scanln(&MobileNo)

					curp.MobileNo = MobileNo

					editPsg(curp)
				case "4":
					fmt.Print("Enter new email address: ")
					fmt.Scanln(&EmailAddress)

					curp.EmailAddress = EmailAddress

					editPsg(curp)
				case "5":
					break editDetailsLoop
				default:
					fmt.Println("Error - Invalid input!")
					continue
				}
			}

		case "4":
			break psgloop
		}
	}
}

func DriverFunctions(curd *Driver) {

	var choice string

drvloop:
	for {
		fmt.Println("\n=======================")
		fmt.Println("|| Driver Functions ||")
		fmt.Println("=======================")
		fmt.Printf("Hello, %s %s! \n", curd.FirstName, curd.LastName)
		fmt.Println("=======================")
		fmt.Println("No assigned rides")
		fmt.Println("=======================\n")
		fmt.Println("1. Edit account information")
		fmt.Println("2. Exit to main menu")
		fmt.Print("Enter an option:")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			var Username string
			var Password string
			var FirstName string
			var LastName string
			var MobileNo string
			var EmailAddress string
			var CarLicenseNo string
		editDetailsLoop:
			for {
				fmt.Println("\n=================")
				fmt.Println("1. Edit username and password")
				fmt.Println("2. Edit first name and last name")
				fmt.Println("3. Edit mobile number")
				fmt.Println("4. Edit email address")
				fmt.Println("5. Edit car license number")
				fmt.Println("6. Exit edit menu")
				fmt.Print("Enter an option:")
				fmt.Scanln(&choice)

				switch choice {
				case "1":
					fmt.Print("Enter new username: ")
					fmt.Scanln(&Username)
					fmt.Print("Enter new password: ")
					fmt.Scanln(&Password)

					curd.Username = Username
					curd.Password = Password

					editDrv(curd)
				case "2":
					fmt.Print("Enter new first name: ")
					fmt.Scanln(&FirstName)
					fmt.Print("Enter new last name: ")
					fmt.Scanln(&LastName)

					curd.FirstName = FirstName
					curd.LastName = LastName

					editDrv(curd)
				case "3":
					fmt.Print("Enter new mobile number: ")
					fmt.Scanln(&MobileNo)

					curd.MobileNo = MobileNo

					editDrv(curd)
				case "4":
					fmt.Print("Enter new email address: ")
					fmt.Scanln(&EmailAddress)

					curd.EmailAddress = EmailAddress

					editDrv(curd)
				case "5":
					fmt.Print("Enter new car license number: ")
					fmt.Scanln(&CarLicenseNo)

					curd.CarLicenseNo = CarLicenseNo

					editDrv(curd)
				case "6":
					break editDetailsLoop
				default:
					fmt.Println("Error - Invalid input!")
					continue
				}
			}

		case "2":
			break drvloop
		}
	}
}

func editPsg(psg *Passenger) {
	url := "http://localhost:6001/api/drive/edit/passenger/" + strconv.Itoa(psg.UserID)
	postBody, _ := json.Marshal(*psg)
	resBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if res.StatusCode == 202 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))
			} else if res.StatusCode == 409 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))

				// This part resets the values of Passenger object CurrentPassenger back to its original due to error
				db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
				if err != nil {
					panic(err.Error())
				}
				defer db.Close()
				var p Passenger
				db.QueryRow("select * from Passengers where Username=? and Password=?", psg.Username, psg.Password).Scan(&p.UserID,
					&p.Username, &p.Password, &p.FirstName, &p.LastName, &p.MobileNo, &p.EmailAddress)
				*psg = p
			}
		}
	}
}

func editDrv(drv *Driver) {
	url := "http://localhost:6001/api/drive/edit/driver/" + strconv.Itoa(drv.UserID)
	postBody, _ := json.Marshal(*drv)
	resBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if res.StatusCode == 202 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))
			} else if res.StatusCode == 409 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))

				// This part resets the values of Driver object CurrentDriver back to its original due to error
				db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/DriveUserDB")
				if err != nil {
					panic(err.Error())
				}
				defer db.Close()
				var d Driver
				db.QueryRow("select * from Passengers where Username=? and Password=?", drv.Username, drv.Password).Scan(&d.UserID,
					&d.Username, &d.Password, &d.FirstName, &d.LastName, &d.MobileNo, &d.EmailAddress)
				*drv = d
			}
		}
	}
}
