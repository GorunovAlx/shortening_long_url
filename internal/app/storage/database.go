package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

/*
type DBStorage struct {
	dsn string

}

func NewDBStorage() *DBStorage {

}

func (dbs *DBStorage) Ping() error {
	return dbs.connection.Ping()
}

func(dbs *DBStorage) ConnectDB() error {
	conn, err := pgx.Connect(context.Background(), dbs.dsn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

}
*/
const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "barkleys"
	dbname   = "dbgolangedu"
)

func TestHandle() error {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	e := conn.Ping(context.Background())

	if e != nil {
		return e
	}

	return nil
}
