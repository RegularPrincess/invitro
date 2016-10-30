// hendlers
package main

import (
	"encoding/json"
	"fmt"
	db "invitro_model"
	"net/http"
	"strconv"
)

var store db.Store

func initHendlers() {
	store = db.GetConnection(19, "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2")
}

func getById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/get/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Fprintf(w, "Uncorrect id")
		return
	}
	analys, err := store.GetById(id)
	analysJson, _ := json.Marshal(analys)
	if err != nil {
		fmt.Fprintf(w, "Not exist id")
		return
	}
	fmt.Fprintln(w, string(analysJson))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}
