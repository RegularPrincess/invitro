// hendlers
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "github.com/RegularPrincess/invitro/invitro_model"
)

var store db.Store

func initHendlers(infoDbConn string) {
	store = db.GetConnection(4, infoDbConn)
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
