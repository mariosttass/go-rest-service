package repository

import (
	"database/sql"

	"github.com/mariosttass/go-rest-service/pkg/model"
)

type Repository struct {
	DB *sql.DB
}

func (db Repository) AddObject(object *model.Object) error {
	var id int
	var createdAt string
	query := `INSERT INTO objects (id, online) VALUES ($1, $2) RETURNING id, created_at`
	err := db.DB.QueryRow(query, object.Number, object.Online).Scan(&id, &createdAt)
	if err != nil {
		return err
	}

	object.ID = id
	object.CreatedAt = createdAt
	return nil
}

func (db Repository) GetObjectById(objectId int) (model.Object, error) {
	item := model.Object{}

	query := `SELECT * FROM objects WHERE id = $1;`
	row := db.DB.QueryRow(query, objectId)
	switch err := row.Scan(&item.ID, &item.Number, &item.Online, &item.CreatedAt); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}

func (db Repository) DeleteObject(objectId int) error {
	query := `DELETE FROM objects WHERE id = $1;`
	_, err := db.DB.Exec(query, objectId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}
