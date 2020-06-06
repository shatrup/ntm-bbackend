package database

// https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
import (
	"github.com/jmoiron/sqlx"
)

var (
	// BhaiFi is the connection handle for the bhaifi database
	BhaiFi *sqlx.DB
	// Radius is the connection handle for the radius database
	//Radius *sqlx.DB
)
