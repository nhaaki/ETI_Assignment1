package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Key struct {
	XMLName xml.Name `xml:"key"`
	Value   string   `xml:"value"`
}

func GetJWT() (string, error) {

	xmlFile, err := os.Open("../key.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var key Key
	xml.Unmarshal(byteValue, &key)

	var mySigningKey = []byte(key.Value)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Krissanawat"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

	fmt.Fprintf(w, string(validToken))
}

func handleRequests() {
	http.HandleFunc("/get/JWT", Index)

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func main() {
	handleRequests()
}
