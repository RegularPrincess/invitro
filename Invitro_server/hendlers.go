// hendlers
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	db "github.com/RegularPrincess/invitro/invitro_model"
)

var store db.Store

func initHendlers(infoDbConn string) {
	store = db.GetConnection(10, infoDbConn)
}

func getById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/get/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Uncorrect id")
		return
	}
	analys, err := store.GetById(id)
	analysJson, _ := json.Marshal(analys)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Not exist id")
		return
	}
	fmt.Fprintln(w, string(analysJson))
}

func index(w http.ResponseWriter, r *http.Request) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]) + "/pages/welcome.html")
	dat, err := ioutil.ReadFile(dir)
	var resp string
	if err != nil {
		log.Println(err)
		resp = "Welcome"
	} else {
		resp = string(dat)
	}
	fmt.Fprintf(w, resp)
}
