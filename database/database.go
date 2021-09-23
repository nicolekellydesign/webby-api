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

CREATE TABLE IF NOT EXISTS gallery_items (
	id TEXT UNIQUE NOT NULL PRIMARY KEY,
	title_line_1 TEXT NOT NULL,
	title_line_2 TEXT NOT NULL,
	thumbnail_location TEXT NOT NULL,
	thumbnail_caption TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS slides (
	id SERIAL PRIMARY KEY,
	gallery_id TEXT NOT NULL,
	name TEXT UNIQUE NOT NULL,
	title TEXT NOT NULL,
	caption TEXT NOT NULL,
	location TEXT NOT NULL,
	CONSTRAINT fk_gallery
		FOREIGN KEY(gallery_id)
		REFERENCES gallery_items(id)
		ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT UNIQUE NOT NULL PRIMARY KEY,
	user_name TEXT UNIQUE NOT NULL,
	user_id SERIAL UNIQUE NOT NULL,
	expires TIMESTAMPTZ NOT NULL,
	CONSTRAINT fk_user_name
		FOREIGN KEY(user_name)
		REFERENCES users(user_name)
		ON DELETE CASCADE,
	CONSTRAINT fk_user_id
		FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);
`

// Connect opens a connection to the database and creates the
// table structure.
func Connect(username, password, database string, log *waterlog.WaterLog) (*DB, error) {
	source := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable timezone=UTC", username, password, database)
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

// AddGalleryItem adds a new gallery item to the database.
func (db DB) AddGalleryItem(item entities.GalleryItem) error {
	tx := db.db.MustBegin()

	sql := `INSERT INTO gallery_items (
		id,
		title_line_1,
		title_line_2,
		thumbnail_location,
		thumbnail_caption
	) VALUES (:id, :title_line_1, :title_line_2, :thumbnail_location, :thumbnail_caption);`

	tx.NamedExec(sql, item)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error adding gallery item to database: %s\n", err)
		return err
	}

	return nil
}

// GetGalleryItems returns all gallery items from the database.
// This makes two requests: one to get the items, and another to
// get all of the slides for each item.
func (db DB) GetGalleryItems() ([]*entities.GalleryItem, error) {
	items := make([]*entities.GalleryItem, 0)

	query := `SELECT
		id,
		title_line_1,
		title_line_2,
		thumbnail_location,
		thumbnail_caption
	FROM gallery_items;`

	if err := db.db.Select(&items, query); err != nil {
		db.log.Errorf("error getting gallery items from the database: %s\n", err)
		return nil, err
	}

	// Get the slides for each gallery item
	for _, item := range items {
		slides := make([]*entities.Slide, 0)
		query := `SELECT
			gallery_id,
			title,
			caption,
			location
		FROM slides WHERE gallery_id=$1;`

		if err := db.db.Select(&slides, query, item.Name); err != nil {
			db.log.Errorf("error getting slide from the database: %s\n", err)
			continue
		}

		item.Slides = slides
	}

	return items, nil
}

// RemoveGalleryItem delets a gallery item from the database.
func (db DB) RemoveGalleryItem(name string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM gallery_items WHERE id=$1;", name)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error removing gallery item from database: %s\n", err)
		return err
	}

	return nil
}

// AddSlide inserts a new slide into the database.
func (db DB) AddSlide(slide entities.Slide) error {
	tx := db.db.MustBegin()

	sql := `INSERT INTO slides (
		gallery_id,
		name,
		title,
		caption,
		location
	) VALUES (
		:gallery_id,
		:name,
		:title,
		:caption,
		:location
	);`

	tx.NamedExec(sql, slide)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error adding slide to database: %s\n", err)
		return err
	}

	return nil
}

// RemoveSlide deletes a slide from the database.
func (db DB) RemoveSlide(galleryID, name string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM slides WHERE gallery_id=$1 AND name=$2;", galleryID, name)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error removing a slide from database: %s\n", err)
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

// GetUser fetches a single user from the database with the
// given ID.
func (db DB) GetUser(id string) (*entities.User, error) {
	var user entities.User
	err := db.db.Get(&user, "SELECT id, user_name FROM users WHERE id=$1;")
	if err != nil {
		db.log.Errorf("error getting user from database: %s\n", err)
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

// RemoveUser deletes a user from the database.
func (db DB) RemoveUser(id string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM users WHERE id=$1;", id)
	if err := tx.Commit(); err != nil {
		db.log.Errorf("error removing user from database: %s\n", err)
		return err
	}

	return nil
}

// GetLogin looks for a matching user for the given username and password. If
// a match is found, a User struct is returned with the id and username.
func (db DB) GetLogin(username, password string) (*entities.User, error) {
	var user entities.User
	err := db.db.Get(&user, "SELECT id, user_name FROM users WHERE user_name=$1 AND (pwdhash = crypt($2, pwdhash));", username, password)
	if err != nil && err != sql.ErrNoRows {
		db.log.Errorf("error checking user login: %s\n", err)
		return nil, err
	}

	return &user, nil
}

// AddSession saves a login session in the database.
func (db DB) AddSession(session *entities.Session) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO sessions VALUES ($1, $2, $3, $4);", session.Token, session.Username, session.ID, session.Expires)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error saving session in database: %s\n", err)
		return err
	}

	return nil
}

// GetSession fetches a session from the database.
func (db DB) GetSession(token string) (*entities.Session, error) {
	var session entities.Session

	sql := `SELECT 
		token,
		user_name,
		user_id,
		expires
	FROM
		sessions
	WHERE
		token = $1;`

	if err := db.db.Get(&session, sql, token); err != nil {
		db.log.Errorf("error getting session from database: %s\n", err)
		return nil, err
	}

	return &session, nil
}

// RemoveSession deletes a session from the database by the session token.
func (db DB) RemoveSession(token string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM sessions WHERE token=$1;", token)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error removing a session from database: %s\n", err)
		return err
	}

	return nil
}

// RemoveSessionForName deletes a session from the database that is
// tied to a particular username.
func (db DB) RemoveSessionForName(username string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM sessions WHERE user_name=$1;", username)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error removing a session from database: %s\n", err)
		return err
	}

	return nil
}

// UpdateSession sets a new expiration time for an existing session.
func (db DB) UpdateSession(session *entities.Session) error {
	tx := db.db.MustBegin()

	sql := `UPDATE sessions
			SET
				expires = :expires
			WHERE
				token = :token;`

	tx.NamedExec(sql, session)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.log.Errorf("error updating a session in the database: %s\n", err)
		return err
	}

	return nil
}
