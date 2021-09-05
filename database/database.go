package database

import (
	"fmt"

	"github.com/DataDrake/waterlog"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nicolekellydesign/webby-api/entities"
)

// DB holds our database connection.
type DB struct {
	db  *sqlx.DB
	log *waterlog.WaterLog
}

var schema = `
CREATE TABLE IF NOT EXISTS auth (
	id SERIAL PRIMARY KEY,
	user_name TEXT NOT NULL,
	pwdhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS photos (
	id SERIAL PRIMARY KEY,
	filename TEXT NOT NULL
);
`

// Connect opens a connection to the database and creates the
// table structure.
func Connect(username, password, database string, log *waterlog.WaterLog) (*DB, error) {
	source := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", username, password, database)
	db, err := sqlx.Connect("pgx", source)
	if err != nil {
		log.Errorf("error connecting to Postgres database: %s\n", err)
		return nil, err
	}
	log.Infof("connected to Postgres database with user '%s' using database name '%s'\n", username, database)

	tx := db.MustBegin()
	tx.MustExec(schema)
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Errorf("error creating table: %s\n", err)
		return nil, err
	}

	self := DB{
		db,
		log,
	}

	return &self, nil
}

// Close closes the database connection and waits for all current
// queries to finish.
func (db DB) Close() {
	if err := db.db.Close(); err != nil {
		db.log.Errorf("error closing database connection: %s\n", err)
	} else {
		db.log.Infoln("closed connection to database")
	}
}

// AddUser inserts a new user with password into the database.
func (db DB) AddUser(user *entities.User) error {
	tx := db.db.MustBegin()
	tx.NamedExec("INSERT INTO auth (user_name, pwdhash) VALUES (:user_name, crypt(:pwdhash, gen_salt('bf')));", user)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error adding user to database: %s\n", err)
		return err
	}

	return nil
}

// CheckLogin tests if the sent user is a valid user in the database.
func (db DB) CheckLogin(user *entities.User) (bool, error) {
	var valid bool
	err := db.db.Get(&valid, "SELECT (pwdhash = crypt($1, pwdhash)) AS pwdhash FROM auth WHERE user_name=$2;", user.Password, user.Username)
	if err != nil {
		db.log.Errorf("error checking user login: %s\n", err)
		return false, err
	}

	return valid, nil
}
