package repository

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	postgresDriver = "postgres"
)

var log logrus.Logger

type DBConnector struct {
	Repository
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
//"postgres://user:passt@host:5432/db"
func ConnString(user, password, host, name string, port int) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
}

func NewDbConnector(connString string) (*DBConnector, error) {
	// creating a dbConnector connection
	log.Info(connString)
	dbConnector, errOpen := sql.Open(postgresDriver, connString)
	if errOpen != nil {
		log.Error("Couldn't connect to DB", errOpen)
		return nil, errOpen
	}

	repository := Repository{DB: dbConnector}
	return &DBConnector{repository}, nil
}

func (c *DBConnector) Stop() {
	c.DB.Close()
}
