package mysql

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type SqlConn struct {
	conn *sql.DB
}

func (c *SqlConn) CreateConnection() {
	db, err := sql.Open("mysql", os.Getenv("db"))

	if err != nil {
		log.Println("mysql conn err: ", err)
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.conn = db
}
