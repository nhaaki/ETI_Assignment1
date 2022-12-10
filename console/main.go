package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// Driver object
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

// Passenger object
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
	// Values that are set each time a user signs in.
	var currentPassenger Passenger
	var currentDriver Driver
	var currentToken string

	// Main loop
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
			// Sign in loop
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

					url := "http://localhost:5000/api/drive/login/passenger"
					loginPayload := map[string]string{
						"Username": Username,
						"Password": Password,
					}
					postBody, _ := json.Marshal(loginPayload)
					resBody := bytes.NewBuffer(postBody)

					if req, err := http.NewRequest("GET", url, resBody); err == nil {
						if res, err := client.Do(req); err == nil {
							defer res.Body.Close()
							if res.StatusCode == 404 {
								fmt.Printf("Error - user not found or incorrect password! \n")
							} else if res.StatusCode == 202 {
								url = "http://localhost:8082/get/JWT"
								if req, err := http.NewRequest("GET", url, nil); err == nil {
									if res, err := client.Do(req); err == nil {
										defer res.Body.Close()
										if res.StatusCode == 409 {
											fmt.Printf("Error in retrieving JWT token... \n")
											continue
										} else if res.StatusCode == 202 {
											body, _ := ioutil.ReadAll(res.Body)
											currentToken = string(body)
										}
									}
								}

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
									// Switch over to passenger menu
									PassengerFunctions(&currentPassenger, &currentToken)
									break signInLoop
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
					url := "http://localhost:5000/api/drive/login/driver"
					loginPayload := map[string]string{
						"Username": Username,
						"Password": Password,
					}
					postBody, _ := json.Marshal(loginPayload)
					resBody := bytes.NewBuffer(postBody)
					if req, err := http.NewRequest("GET", url, resBody); err == nil {
						if res, err := client.Do(req); err == nil {
							defer res.Body.Close()
							if res.StatusCode == 404 {
								fmt.Printf("Error - user not found or incorrect password! \n\n")
							} else if res.StatusCode == 202 {

								url = "http://localhost:8082/get/JWT"
								if req, err := http.NewRequest("GET", url, nil); err == nil {
									if res, err := client.Do(req); err == nil {
										defer res.Body.Close()
										if res.StatusCode == 409 {
											fmt.Printf("Error in retrieving JWT token... \n")
											continue
										} else if res.StatusCode == 202 {
											body, _ := ioutil.ReadAll(res.Body)
											currentToken = string(body)
										}
									}
								}

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
									// Switch over to driver menu
									DriverFunctions(&currentDriver, &currentToken)
									break signInLoop
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
			// Account creation loop
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

						client := &http.Client{}
						url := "http://localhost:5000/api/drive/checkusername?username=" + Username + "&type=Passenger"
						if req, err := http.NewRequest("POST", url, nil); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 404 || res.StatusCode == 409 {
									body, _ := ioutil.ReadAll(res.Body)
									fmt.Println(string(body))
									continue
								}
							}
						}

						fmt.Print("Enter your password: ")
						fmt.Scanln(&Password)
						fmt.Print("Enter your first name: ")
						fmt.Scanln(&FirstName)
						fmt.Print("Enter your last name: ")
						fmt.Scanln(&LastName)
						fmt.Print("Enter your mobile no (8 characters): ")
						fmt.Scanln(&MobileNo)
						fmt.Print("Enter your email address: ")
						fmt.Scanln(&EmailAddress)

						newAccount := Passenger{999, Username, Password, FirstName, LastName, MobileNo, EmailAddress}
						url = "http://localhost:5000/api/drive/create/passenger"
						postBody, _ := json.Marshal(newAccount)
						resBody := bytes.NewBuffer(postBody)
						client = &http.Client{}
						if req, err := http.NewRequest("POST", url, resBody); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 202 {
									fmt.Printf("Account %s created successfully!\n", newAccount.Username)
									break psgloop

								} else if res.StatusCode == 409 {
									body, _ := ioutil.ReadAll(res.Body)
									fmt.Println(string(body))
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

						client := &http.Client{}
						url := "http://localhost:5000/api/drive/checkusername?username=" + Username + "&type=Driver"
						if req, err := http.NewRequest("POST", url, nil); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 404 || res.StatusCode == 409 {
									body, _ := ioutil.ReadAll(res.Body)
									fmt.Println(string(body))
									continue
								}
							}
						}

						fmt.Print("Enter your password: ")
						fmt.Scanln(&Password)
						fmt.Print("Enter your first name: ")
						fmt.Scanln(&FirstName)
						fmt.Print("Enter your last name: ")
						fmt.Scanln(&LastName)
						fmt.Print("Enter your mobile no (8 characters): ")
						fmt.Scanln(&MobileNo)
						fmt.Print("Enter your email address: ")
						fmt.Scanln(&EmailAddress)
						fmt.Print("Enter your ID number: ")
						fmt.Scanln(&IdNo)
						fmt.Print("Enter your car license number: ")
						fmt.Scanln(&CarLicenseNo)

						newAccount := Driver{999, Username, Password, FirstName, LastName, MobileNo, EmailAddress, IdNo, CarLicenseNo}
						url = "http://localhost:5000/api/drive/create/driver"
						postBody, _ := json.Marshal(newAccount)
						resBody := bytes.NewBuffer(postBody)
						client = &http.Client{}
						if req, err := http.NewRequest("POST", url, resBody); err == nil {
							if res, err := client.Do(req); err == nil {
								defer res.Body.Close()
								if res.StatusCode == 202 {
									fmt.Printf("Account %s created successfully!\n", newAccount.Username)
									break drvloop
								} else if res.StatusCode == 409 {
									body, _ := ioutil.ReadAll(res.Body)
									fmt.Println(string(body))
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

// Passenger menu and functions
func PassengerFunctions(curp *Passenger, currentToken *string) {
	var choice string

	// Passenger menu loop
psgloop:
	for {
		fmt.Println("\n=========================")
		fmt.Println("|| Passenger Functions ||")
		fmt.Println("=========================")
		fmt.Printf("Hello, %s %s! \n", curp.FirstName, curp.LastName)
		fmt.Println("=======================")
		fmt.Println("1. Book a ride")
		fmt.Println("2. Retrieve all trips")
		fmt.Println("3. Edit account information")
		fmt.Println("4. Exit to main menu")
		fmt.Print("Enter an option:")
		fmt.Scanln(&choice)

		switch choice {
		// Book a ride
		case "1":
			var pcPickup string
			var pcDropoff string

			fmt.Println("\n=================")
			fmt.Print("Enter postal code of pick-up location: ")
			fmt.Scanln(&pcPickup)
			fmt.Print("Enter postal code of drop-off location: ")
			fmt.Scanln(&pcDropoff)

			client := &http.Client{}
			url := "http://localhost:6002/api/drive/passenger/assigndriver/" + strconv.Itoa(curp.UserID)
			loginPayload := map[string]string{
				"pcPickup":  pcPickup,
				"pcDropoff": pcDropoff,
			}
			postBody, _ := json.Marshal(loginPayload)
			resBody := bytes.NewBuffer(postBody)

			if req, err := http.NewRequest("POST", url, resBody); err == nil {
				req.Header.Set("Token", *currentToken)
				if res, err := client.Do(req); err == nil {
					defer res.Body.Close()
					if res.StatusCode == 404 {
						body, _ := ioutil.ReadAll(res.Body)
						fmt.Println(string(body))
					} else if res.StatusCode == 202 {
						// If driver is successfully assigned to passenger, continue
						body, _ := ioutil.ReadAll(res.Body)
						fmt.Println(string(body))
						fmt.Println("Waiting for driver to start ride...")
						fmt.Println("+====+")
						fmt.Println("|(::)|")
						fmt.Println("| )( |")
						fmt.Println("|(..)|")
						fmt.Println("+====+")

						// Check the ride status. Once the ride is ended, value would be 0. Otherwise, 1.
						var count int
						url := "http://localhost:6002/api/drive/checkactiveride?userid=" + strconv.Itoa(curp.UserID)
						if req, err := http.NewRequest("GET", url, nil); err == nil {
							if res, err := client.Do(req); err == nil {
								body, _ := ioutil.ReadAll(res.Body)
								err := json.Unmarshal(body, &count)
								if err != nil {
									panic(err.Error())
								}
							}
						}

						// Start a loop that keeps calling this method through the API. This method is explained in detail
						// in /bookride/main.go
						for count != 0 {
							url := "http://localhost:6002/api/drive/driver/ridefunctions/" + strconv.Itoa(curp.UserID)
							if req, err := http.NewRequest("POST", url, nil); err == nil {
								req.Header.Set("Token", *currentToken)
								if res, err := client.Do(req); err == nil {
									defer res.Body.Close()
									if res.StatusCode == 201 {
										body, _ := ioutil.ReadAll(res.Body)
										fmt.Println(string(body))
									} else if res.StatusCode == 202 {
										body, _ := ioutil.ReadAll(res.Body)
										fmt.Println(string(body))

									} else if res.StatusCode == 200 {
										body, _ := ioutil.ReadAll(res.Body)
										fmt.Println(string(body))
									}
								}
							}

							// After receiving the response from the API, check if count value has changed.
							// If ride ended, value should change to 0, thus ending the loop.
							url = "http://localhost:6002/api/drive/checkactiveride?userid=" + strconv.Itoa(curp.UserID)
							if req, err := http.NewRequest("GET", url, nil); err == nil {
								if res, err := client.Do(req); err == nil {
									body, _ := ioutil.ReadAll(res.Body)
									err := json.Unmarshal(body, &count)
									if err != nil {
										panic(err.Error())
									}
								}
							}
						}
					}
				}
			}
		// Get ride history in reverse chronological order.
		case "2":
			client := &http.Client{}
			url := "http://localhost:6005/api/drive/passenger/gethistory/" + strconv.Itoa(curp.UserID)
			if req, err := http.NewRequest("GET", url, nil); err == nil {
				req.Header.Set("Token", *currentToken)
				if res, err := client.Do(req); err == nil {
					defer res.Body.Close()
					if res.StatusCode == 202 {
						body, _ := ioutil.ReadAll(res.Body)
						var allRides map[int]map[string]string
						err := json.Unmarshal(body, &allRides)
						if err != nil {
							panic(err.Error())
						}

						for index, details := range allRides {
							fmt.Printf("\n(%v)\n", index)
							fmt.Println("=================")
							fmt.Println("Ride Date: " + details["rideDate"])
							fmt.Println("\nDriver name: " + details["driverFirstName"] + " " + details["driverLastName"])
							fmt.Println("Driver Car License No. : " + details["carLicenseNo"])
							fmt.Println("\nPick-up postal code: " + details["pcPickup"])
							fmt.Println("Drop-off postal code: " + details["pcDropOff"])

							fmt.Println("=================")
						}

					} else if res.StatusCode == 404 {
						body, _ := ioutil.ReadAll(res.Body)
						fmt.Println("\n" + string(body))
					}
				}
			}
		// Edit passenger details.
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

					editPsg(curp, *currentToken)
				case "2":
					fmt.Print("Enter new first name: ")
					fmt.Scanln(&FirstName)
					fmt.Print("Enter new last name: ")
					fmt.Scanln(&LastName)

					curp.FirstName = FirstName
					curp.LastName = LastName

					editPsg(curp, *currentToken)
				case "3":
					fmt.Print("Enter new mobile number: ")
					fmt.Scanln(&MobileNo)

					curp.MobileNo = MobileNo

					editPsg(curp, *currentToken)
				case "4":
					fmt.Print("Enter new email address: ")
					fmt.Scanln(&EmailAddress)

					curp.EmailAddress = EmailAddress

					editPsg(curp, *currentToken)
				case "5":
					break editDetailsLoop
				default:
					fmt.Println("Error - Invalid input!")
					continue
				}
			}
		// Log out. Void out the values set in the beginning.
		case "4":
			*curp = Passenger{}
			*currentToken = ""
			break psgloop
		}
	}
}

// Driver menu and functions
func DriverFunctions(curd *Driver, currentToken *string) {
	var choice string
	// Driver menu loop
drvloop:
	for {
		fmt.Println("\n=======================")
		fmt.Println("|| Driver Functions ||")
		fmt.Println("=======================")
		fmt.Printf("Hello, %s %s! \n", curd.FirstName, curd.LastName)
		fmt.Println("=======================")
		// In order to print the status of the driver, two things need to be done.
		// 1. Get the status from the database,
		// 2. Use the printDriverStatus() method in order to print out the extra details such as the passenger being assigned.
		var status string
		client := &http.Client{}
		url := "http://localhost:6002/api/drive/getdriverstatus?userid=" + strconv.Itoa(curd.UserID)
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 202 {
					body, _ := ioutil.ReadAll(res.Body)
					status = string(body)
				}
			}
		}
		fmt.Println("Status: " + status)

		currState := printDriverStatus(status, curd, *currentToken)

		fmt.Println("=======================")
		fmt.Println("\n1. Edit account information")
		fmt.Println("2. Refresh page")
		fmt.Println("3. Exit to main menu")
		fmt.Print("Enter an option:")
		fmt.Scanln(&choice)

		switch choice {
		// Edit account information
		case "1":
			var Username string
			var Password string
			var FirstName string
			var LastName string
			var MobileNo string
			var EmailAddress string
			var CarLicenseNo string
			// Edit account loop
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

					editDrv(curd, *currentToken)
				case "2":
					fmt.Print("Enter new first name: ")
					fmt.Scanln(&FirstName)
					fmt.Print("Enter new last name: ")
					fmt.Scanln(&LastName)

					curd.FirstName = FirstName
					curd.LastName = LastName

					editDrv(curd, *currentToken)
				case "3":
					fmt.Print("Enter new mobile number: ")
					fmt.Scanln(&MobileNo)

					curd.MobileNo = MobileNo

					editDrv(curd, *currentToken)
				case "4":
					fmt.Print("Enter new email address: ")
					fmt.Scanln(&EmailAddress)

					curd.EmailAddress = EmailAddress

					editDrv(curd, *currentToken)
				case "5":
					fmt.Print("Enter new car license number: ")
					fmt.Scanln(&CarLicenseNo)

					curd.CarLicenseNo = CarLicenseNo

					editDrv(curd, *currentToken)
				case "6":
					break editDetailsLoop
				default:
					fmt.Println("Error - Invalid input!")
					continue
				}
			}
		// Refresh the page, usually used to check if the driver has been assigned a passenger.
		case "2":
			continue
		// Start a ride that has been assigned.
		case "a":
			if currState != 1 {
				fmt.Println("Error - Invalid input!")
			} else {
				client := &http.Client{}
				url := "http://localhost:6003/api/drive/startride/" + strconv.Itoa(curd.UserID)
				if req, err := http.NewRequest("POST", url, nil); err == nil {
					req.Header.Set("Token", *currentToken)
					_, err := client.Do(req)
					if err != nil {
						panic(err.Error())
					}
				}

			}
		// End the ride that has been assigned.
		case "b":
			if currState != 2 {
				fmt.Println("Error - Invalid input!")
			} else {
				client := &http.Client{}
				url := "http://localhost:6003/api/drive/endride/" + strconv.Itoa(curd.UserID)
				if req, err := http.NewRequest("POST", url, nil); err == nil {
					req.Header.Set("Token", *currentToken)
					if res, err := client.Do(req); err == nil {
						defer res.Body.Close()
						if res.StatusCode == 202 {
							body, _ := ioutil.ReadAll(res.Body)
							resBody := bytes.NewBuffer(body)
							url = "http://localhost:6005/api/drive/addtohistory"
							if req, err := http.NewRequest("POST", url, resBody); err == nil {
								req.Header.Set("Token", *currentToken)
								if res, err := client.Do(req); err == nil {
									defer res.Body.Close()
									if res.StatusCode == 409 {
										body, _ := ioutil.ReadAll(res.Body)
										fmt.Println(string(body))
									}
								}
							}
						}
					}
				}
			}
		// Cancel a ride that has been assigned.
		case "c":
			if currState != 1 {
				fmt.Println("Error - Invalid input!")
			} else {
				client := &http.Client{}
				url := "http://localhost:6003/api/drive/cancelride/" + strconv.Itoa(curd.UserID)
				if req, err := http.NewRequest("POST", url, nil); err == nil {
					req.Header.Set("Token", *currentToken)
					if res, err := client.Do(req); err == nil {
						defer res.Body.Close()
						if res.StatusCode == 202 {
							body, _ := ioutil.ReadAll(res.Body)
							fmt.Println(string(body))
						}
					}
				}
			}
		// Log out. Need to check if driver is in the middle of a trip or just been assigned, as they can't cancel if
		// the trip has already started. Remove the driver's record in live trips table, and void the values set earlier.
		case "3":
			client := &http.Client{}
			url := "http://localhost:6002/api/drive/getdriverstatus" + strconv.Itoa(curd.UserID)
			if req, err := http.NewRequest("POST", url, nil); err == nil {
				if res, err := client.Do(req); err == nil {
					defer res.Body.Close()
					if res.StatusCode == 202 {
						body, _ := ioutil.ReadAll(res.Body)
						status = string(body)
					}
				}
			}

			if status == "Ongoing" {
				fmt.Println("You can't log out while driving a passenger!")
				continue
			} else if status == "Assigned" {
				client := &http.Client{}
				url := "http://localhost:6003/api/drive/cancelride/" + strconv.Itoa(curd.UserID)
				if req, err := http.NewRequest("POST", url, nil); err == nil {
					req.Header.Set("Token", *currentToken)
					if res, err := client.Do(req); err == nil {
						defer res.Body.Close()
						if res.StatusCode == 202 {
							body, _ := ioutil.ReadAll(res.Body)
							fmt.Println(string(body))
						}
					}
				}
			}
			url = "http://localhost:6002/api/drive/driver/setstatus?userid=" + strconv.Itoa(curd.UserID)
			if req, err := http.NewRequest("DELETE", url, nil); err == nil {
				req.Header.Set("Token", *currentToken)
				if res, err := client.Do(req); err == nil {
					defer res.Body.Close()
					if res.StatusCode == 202 {
						body, _ := ioutil.ReadAll(res.Body)
						fmt.Println(string(body))
					}
				}
			}
			*curd = Driver{}
			*currentToken = ""
			break drvloop

		default:
			fmt.Println("Error - Invalid input!")
		}
	}
}

// Function used to edit passenger details.
func editPsg(psg *Passenger, currentToken string) {
	url := "http://localhost:5000/api/drive/edit/passenger/" + strconv.Itoa(psg.UserID)
	postBody, _ := json.Marshal(*psg)
	resBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		req.Header.Set("Token", currentToken)
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if res.StatusCode == 202 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))
			} else if res.StatusCode == 409 {
				body, _ := ioutil.ReadAll(res.Body)
				var unpwPayload map[string]string
				err := json.Unmarshal(body, &unpwPayload)
				if err != nil {
					panic(err.Error())
				}
				// If there was an error in changing the username and password (username has to be unique),
				// the API returns the initial username and password to be assigned back to the current Passenger object.
				psg.Username = unpwPayload["Username"]
				psg.Password = unpwPayload["Password"]
				fmt.Println("Error - Username taken.")
			}
		}
	}
}

// Function used to edit driver details.
func editDrv(drv *Driver, currentToken string) {
	url := "http://localhost:6001/api/drive/edit/driver/" + strconv.Itoa(drv.UserID)
	postBody, _ := json.Marshal(*drv)
	resBody := bytes.NewBuffer(postBody)
	client := &http.Client{}
	if req, err := http.NewRequest("PUT", url, resBody); err == nil {
		req.Header.Set("Token", currentToken)
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			if res.StatusCode == 202 {
				body, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body))
			} else if res.StatusCode == 409 {
				body, _ := ioutil.ReadAll(res.Body)
				var unpwPayload map[string]string
				err := json.Unmarshal(body, &unpwPayload)
				if err != nil {
					panic(err.Error())
				}
				// If there was an error in changing the username and password (username has to be unique),
				// the API returns the initial username and password to be assigned back to the current Driver object.
				drv.Username = unpwPayload["Username"]
				drv.Password = unpwPayload["Password"]
				fmt.Println("Error - Username taken.")
			}
		}
	}
}

// Print different messages and details according to driver's current status
func printDriverStatus(status string, curd *Driver, currentToken string) (currState int) {

	switch status {
	case "Assigned":

		var rideDetails map[string]string
		client := &http.Client{}
		url := "http://localhost:6002/api/drive/getridedetails?userid=" + strconv.Itoa(curd.UserID)
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			req.Header.Set("Token", currentToken)
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 202 {
					body, _ := ioutil.ReadAll(res.Body)
					err := json.Unmarshal(body, &rideDetails)
					if err != nil {
						panic(err.Error())
					}
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(body))
				}
			}
		}

		fmt.Printf("Assigned to %s %s\n", rideDetails["pFirstName"], rideDetails["pLastName"])
		fmt.Printf("Pick-up postal code: %s\n", rideDetails["pcPickUp"])
		fmt.Printf("Drop-off postal code: %s\n", rideDetails["pcDropOff"])
		fmt.Println("a. Start ride")
		fmt.Println("c. Cancel ride")
		return 1

	case "Available":
		fmt.Println("No assigned rides")
		return 0

	case "Ongoing":
		var rideDetails map[string]string
		client := &http.Client{}
		url := "http://localhost:6002/api/drive/getridedetails?userid=" + strconv.Itoa(curd.UserID)
		if req, err := http.NewRequest("GET", url, nil); err == nil {
			req.Header.Set("Token", currentToken)
			if res, err := client.Do(req); err == nil {
				defer res.Body.Close()
				if res.StatusCode == 202 {
					body, _ := ioutil.ReadAll(res.Body)
					err := json.Unmarshal(body, &rideDetails)
					if err != nil {
						panic(err.Error())
					}
				} else {
					body, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(body))
				}
			}
		}

		fmt.Printf("Assigned to %s %s\n", rideDetails["pFirstName"], rideDetails["pLastName"])
		fmt.Printf("Pick-up postal code: %s\n", rideDetails["pcPickUp"])
		fmt.Printf("Drop-off postal code: %s\n", rideDetails["pcDropOff"])
		fmt.Println("b. Stop ride")
		return 2
	}

	return 0
}
