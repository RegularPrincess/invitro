// invitro_server project main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RegularPrincess/invitro/invitro_parser"
)

func main() {
	link := "/analizes/for-doctors/"
	dbConnInf := "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2"

	parser := invitro_parser.GetParser(dbConnInf)
	if parser.NeedParse() {
		fmt.Println("scraper start")
		go parser.Scrape(link)
	}

	initHendlers(dbConnInf)
	http.HandleFunc("/get/", getById)
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
