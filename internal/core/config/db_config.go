package config

import (
	"strings"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

/*
Generate the connection string for the database
Returns the connection string
*/
func (db *Database) ConnectionString() string {
	var b strings.Builder

	b.WriteString("postgres://")
	b.WriteString(db.User)
	b.WriteString(":")
	b.WriteString(db.Password)
	b.WriteString("@")
	b.WriteString(db.Host)
	b.WriteString(":")
	b.WriteString(db.Port)
	b.WriteString("/")
	b.WriteString(db.Database)
	b.WriteString("?sslmode=")
	b.WriteString(db.SSLMode)

	return b.String()
}
