// db
package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db sql.DB
}

func GetConnection(countConns int, connInfo string) Store {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(countConns)
	return Store{*db}
}

func (this *Store) GetById(id int) (*Analysis, error) {
	result := store.db.QueryRow("SELECT * FROM analysis WHERE (id = $1)", id)
	analys := new(Analysis)
	err := result.Scan(&analys.Id, &analys.Kind, &analys.Subtype, &analys.Name, &analys.Description)
	if err != nil {
		return nil, err
	}
	return analys, nil
}

func (this *Store) addAnalysis(a *Analysis) error {
	_, err := this.db.Exec("INSERT INTO analysis (type, subtype, name, description) VALUES ($1, $2, $3, $4);",
		a.Kind, a.Subtype, a.Name, a.Description)
	if err != nil {
		return err
	}
	return nil
}
