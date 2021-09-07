package database

import (
	"database/sql"
	"fmt"

	"github.com/DataDrake/waterlog"
	// This is commented because I guess that's all sqlx needs
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
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	user_name TEXT UNIQUE NOT NULL,
	pwdhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS photos (
	id SERIAL PRIMARY KEY,
	file_name TEXT NOT NULL
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

// AddPhoto inserts a new photo into the database.
func (db DB) AddPhoto(fileName string) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO photos (file_name) VALUES ($1);", fileName)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error adding photo to database: %s\n", err)
		return err
	}

	return nil
}

// GetPhotos fetches all photos from the database.
func (db DB) GetPhotos() ([]*entities.Photo, error) {
	ret := make([]*entities.Photo, 0)
	if err := db.db.Select(&ret, "SELECT file_name FROM photos;"); err != nil {
		db.log.Errorf("error getting photos from database: %s\n", err)
		return nil, err
	}

	return ret, nil
}

// RemovePhoto deletes a photo with the given file name from
// the database.
func (db DB) RemovePhoto(fileName string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM photos WHERE file_name=$1;", fileName)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error removing photo from database: %s\n", err)
		return err
	}

	return nil
}

// AddUser inserts a new user into the database.
func (db DB) AddUser(username, password string) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO users (user_name, pwdhash) VALUES ($1, crypt($2, gen_salt('bf')));", username, password)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error adding user to database: %s\n", err)
		return err
	}

	return nil
}

// GetLogin looks for a matching user for the given username and password. If
// a match is found, a User struct is returned with the id and username.
func (db DB) GetLogin(username, password string) (*entities.User, error) {
	var user entities.User
	err := db.db.Get(&user, "SELECT id, user_name FROM users WHERE user_name=$1 AND (pwdhash = crypt($2, pwdhash));", user.Username, user.Password)
	if err != nil && err != sql.ErrNoRows {
		db.log.Errorf("error checking user login: %s\n", err)
		return nil, err
	}

	return &user, nil
}

// GetUsers fetches all of the users from the database.
func (db DB) GetUsers() ([]*entities.User, error) {
	ret := make([]*entities.User, 0)
	if err := db.db.Select(&ret, "SELECT id, user_name FROM users;"); err != nil {
		db.log.Errorf("error getting users from database: %s\n", err)
		return nil, err
	}

	return ret, nil
}
