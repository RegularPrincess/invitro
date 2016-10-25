package invitro_server

import (
	//"fmt"
	//"html"
	"log"
	"net/http"
)

func main() {
	//	store = GetConnection(19, "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2")
	//	result := store.db.QueryRow("SELECT * FROM analysis WHERE (id = $1)", 1000)
	//	analys := new(Analysis)
	//	err := result.Scan(&analys.Id, &analys.Kind, &analys.Subtype, &analys.Name, &analys.Description)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(analys.Kind)
	initHendlers()
	http.HandleFunc("/get/", getById)
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
