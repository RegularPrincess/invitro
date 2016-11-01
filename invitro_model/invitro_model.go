// invitro_model project invitro_model.go
package invitro_model

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Analysis struct {
	Id          int    `json:"id"`
	Kind        string `json:"kind"`
	Subtype     string `json:"subtype"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

func (this *Store) CloseConn() {
	this.db.Close()
}

//Проверяет заполненность базы
func (this *Store) DBFilled() (bool, error) {
	result := this.db.QueryRow("SELECT COUNT(id) FROM analysis;")
	var count int
	err := result.Scan(&count)
	filled := count > 1000
	if err != nil {
		return false, err
	}
	return filled, nil
}

//Очищает базу и обнуляет счетчик id
func (this *Store) Clean() error {
	_, err := this.db.Exec("DELETE FROM analysis;")
	if err != nil {
		return err
	}
	_, err = this.db.Exec("ALTER SEQUENCE analysis_id_seq RESTART WITH 1;")
	if err != nil {
		return err
	}
	return nil
}

func (this *Store) GetById(id int) (*Analysis, error) {
	result := this.db.QueryRow("SELECT * FROM analysis WHERE (id = $1)", id)
	analys := new(Analysis)
	err := result.Scan(&analys.Id, &analys.Kind, &analys.Subtype, &analys.Name, &analys.Description)
	if err != nil {
		return nil, err
	}
	return analys, nil
}

func (this *Store) AddAnalysis(a *Analysis) error {
	_, err := this.db.Exec("INSERT INTO analysis (type, subtype, name, description) VALUES ($1, $2, $3, $4);",
		a.Kind, a.Subtype, a.Name, a.Description)
	if err != nil {
		return err
	}
	return nil
}
