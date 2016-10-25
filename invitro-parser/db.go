// db
package invitro_parser

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	db sql.DB
}

func GetConnection(countConns int, connInfo string) Store {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(countConns)
	return Store{*db}
}

func (this *Store) addAnalysis(a *Analysis) {
	_, err := this.db.Exec("INSERT INTO analysis (type, subtype, name, description) VALUES ($1, $2, $3, $4);",
		a.kind, a.subtype, a.name, a.description)
	if err != nil {
		log.Println(err)
	}
}
