// Invitro_server project main.go
package main

import (
	//parser "invitro_parser"
	"log"
	"net/http"
)

func main() {
	//	link := "/analizes/for-doctors/"
	//	dbConnInf := "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2"
	//	parser.Scrape(link, dbConnInf)

	initHendlers()
	http.HandleFunc("/get/", getById)
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

//	link := "/analizes/for-doctors/"
//	dbConnInf := "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2"
//	Scrape(link, dbConnInf)
