package database

import (
	"database/sql"
	"fmt"

	// This is commented because I guess that's all sqlx needs
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/nicolekellydesign/webby-api/entities"
)

// DB holds our database connection.
type DB struct {
	db *sqlx.DB
}

var tables = []string{
	`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		user_name TEXT UNIQUE NOT NULL,
		pwdhash TEXT NOT NULL,
		protected BOOL NOT NULL
	);`,

	`CREATE TABLE photos (
		id SERIAL PRIMARY KEY,
		file_name TEXT NOT NULL
	);`,

	`CREATE TABLE gallery_items (
		id TEXT UNIQUE NOT NULL PRIMARY KEY,
		title TEXT NOT NULL,
		caption TEXT NOT NULL,
		project_info TEXT NOT NULL,
		thumbnail TEXT NOT NULL,
		embed_url TEXT
	);`,

	`CREATE TABLE project_images (
		id SERIAL PRIMARY KEY,
		gallery_id TEXT NOT NULL,
		file_name VARCHAR(255) NOT NULL,
		CONSTRAINT fk_gallery
			FOREIGN KEY(gallery_id)
			REFERENCES gallery_items(id)
			ON DELETE CASCADE
	);`,

	`CREATE TABLE sessions (
		token TEXT UNIQUE NOT NULL PRIMARY KEY,
		user_name TEXT UNIQUE NOT NULL,
		user_id SERIAL UNIQUE NOT NULL,
		created TIMESTAMPTZ NOT NULL,
		max_age BIGINT NOT NULL,
		CONSTRAINT fk_user_name
			FOREIGN KEY(user_name)
			REFERENCES users(user_name)
			ON DELETE CASCADE,
		CONSTRAINT fk_user_id
			FOREIGN KEY(user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
	);`,
}

// Connect opens a connection to the database and creates the
// table structure.
func Connect(username, password, database string) (*DB, error) {
	source := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable timezone=UTC", username, password, database)
	db, err := sqlx.Connect("pgx", source)
	if err != nil {
		return nil, err
	}

	self := DB{
		db,
	}

	return &self, nil
}

// Close closes the database connection and waits for all current
// queries to finish.
func (db DB) Close() {
	db.db.Close()
}

// InitSchema creates our tables in the database.
func (db DB) InitSchema() error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	for _, table := range tables {
		tx.Exec(table)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// AddPhoto inserts a new photo into the database.
func (db DB) AddPhoto(fileName string) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO photos (file_name) VALUES ($1);", fileName)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetPhotos fetches all photos from the database.
func (db DB) GetPhotos() ([]*entities.Photo, error) {
	ret := make([]*entities.Photo, 0)
	if err := db.db.Select(&ret, "SELECT file_name FROM photos;"); err != nil {
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
		return err
	}

	return nil
}

// AddGalleryItem adds a new gallery item to the database.
func (db DB) AddGalleryItem(item entities.GalleryItem) error {
	tx := db.db.MustBegin()

	sql := `INSERT INTO gallery_items (
		id,
		title,
		caption,
		project_info,
		thumbnail,
		embed_url
	) VALUES ($1, $2, $3, $4, $5, $6);`

	tx.MustExec(sql, item.Name, item.Title, item.Caption, item.ProjectInfo, item.Thumbnail, item.EmbedURL.String)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// ChangeProjectThumbnail sets a new thumbnail for a project.
func (db DB) ChangeProjectThumbnail(name, newThumb string) error {
	tx := db.db.MustBegin()
	tx.MustExec("UPDATE gallery_items SET thumbnail=$1 WHERE id=$2;", newThumb, name)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetProject retrieves a project from the database with the given name.
func (db DB) GetProject(name string) (*entities.GalleryItem, error) {
	var project entities.GalleryItem

	query := "SELECT title, caption, project_info, thumbnail, embed_url FROM gallery_items WHERE id=$1;"
	if err := db.db.Get(&project, query, name); err != nil {
		return nil, err
	}

	project.Name = name

	images := make([]string, 0)
	query = "SELECT file_name FROM project_images WHERE gallery_id=$1;"

	if err := db.db.Select(&images, query, name); err != nil {
		return nil, err
	}

	project.Images = images
	return &project, nil
}

// UpdateProject sets the title, caption, project info, and embed URL fields for a project
// with the same name in the database.
func (db DB) UpdateProject(project *entities.GalleryItem) error {
	tx := db.db.MustBegin()

	sql := `
	UPDATE
		gallery_items
	SET
		title = $1,
		caption = $2,
		project_info = $3,
		embed_url = $4
	WHERE
		id = $5;
	`

	tx.MustExec(sql, project.Title, project.Caption, project.ProjectInfo, project.EmbedURL.String, project.Name)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetGalleryItems returns all gallery items from the database.
func (db DB) GetGalleryItems() ([]*entities.GalleryItem, error) {
	items := make([]*entities.GalleryItem, 0)

	query := `SELECT
		id,
		title,
		caption,
		project_info,
		thumbnail,
		embed_url
	FROM gallery_items;`

	if err := db.db.Select(&items, query); err != nil {
		return nil, err
	}

	// Get the project images for each gallery item
	for _, item := range items {
		images := make([]string, 0)
		query := `SELECT
			file_name
		FROM project_images WHERE gallery_id=$1;`

		if err := db.db.Select(&images, query, item.Name); err != nil {
			continue
		}

		item.Images = images
	}

	return items, nil
}

// RemoveGalleryItem delets a gallery item from the database.
func (db DB) RemoveGalleryItem(name string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM gallery_items WHERE id=$1;", name)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Add ProjectImage inserts a file name for an image in a project to the database.
func (db DB) AddProjectImage(galleryID, filename string) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO project_images (gallery_id, file_name) VALUES ($1, $2);", galleryID, filename)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// RemoveProjectImage deletes a project image entry from the database.
func (db DB) RemoveProjectImage(galleryID, filename string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM project_images WHERE gallery_id=$1 AND file_name=$2;", galleryID, filename)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// AddUser inserts a new user into the database.
func (db DB) AddUser(username, password string, protected bool) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO users (user_name, pwdhash, protected) VALUES ($1, crypt($2, gen_salt('bf')), $3);", username, password, protected)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetUser fetches a single user from the database with the
// given ID.
func (db DB) GetUser(id string) (*entities.User, error) {
	var user entities.User
	err := db.db.Get(&user, "SELECT id, user_name, protected FROM users WHERE id=$1;", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUsers fetches all of the users from the database.
func (db DB) GetUsers() ([]*entities.User, error) {
	ret := make([]*entities.User, 0)
	if err := db.db.Select(&ret, "SELECT id, user_name, protected FROM users;"); err != nil {
		return nil, err
	}

	return ret, nil
}

// RemoveUser deletes a user from the database.
func (db DB) RemoveUser(id string) error {
	tx := db.db.MustBegin()
	tx.MustExec("DELETE FROM users WHERE id=$1;", id)
	if err := tx.Commit(); err != nil {
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
		return nil, err
	}

	return &user, nil
}

// AddSession saves a login session in the database.
func (db DB) AddSession(session *entities.Session) error {
	tx := db.db.MustBegin()
	tx.MustExec("INSERT INTO sessions VALUES ($1, $2, $3, $4, $5);", session.Token, session.Username, session.ID, session.Created, session.MaxAge)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetSession fetches a session from the database.
func (db DB) GetSession(token string) (*entities.Session, error) {
	var session entities.Session

	query := `SELECT 
		token,
		user_name,
		user_id,
		created,
		max_age
	FROM
		sessions
	WHERE
		token = $1;`

	err := db.db.Get(&session, query, token)
	if err != nil && err != sql.ErrNoRows {
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
		return err
	}

	return nil
}

// UpdateSession sets a new expiration time for an existing session.
func (db DB) UpdateSession(session *entities.Session) error {
	tx := db.db.MustBegin()

	sql := `UPDATE sessions
			SET
				max_age = :max_age
			WHERE
				token = :token;`

	tx.NamedExec(sql, session)

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
