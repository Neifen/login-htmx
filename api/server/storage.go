package server

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(u *UserType) error
	DeleteUser(u *UserType) error
	ReadUserByEmail(email string) (*UserType, error)
	UpdateUser(u *UserType) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PW")

	// connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=verify-full", host, port, user, password, dbname)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New("error 130: could not initialize db")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.New("error 131: could not initialize db")
	}

	_, err = db.Query("CREATE EXTENSION IF NOT EXISTS citext;")
	if err != nil {
		return nil, errors.New("error 132: could not initialize db")
	}

	_, err = db.Query(`
	CREATE TABLE 
	IF NOT EXISTS 
	users(
	id SERIAL PRIMARY KEY,
	name varchar not null,
	email citext not null unique,
	pw bytea not null,
	uid varchar not null,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
	)
	`)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error 133: could not initialize db")
	}

	_, err = db.Query("CREATE UNIQUE INDEX IF NOT EXISTS users_unique_lower_email_idx ON users (lower(email));")
	if err != nil {
		return nil, errors.New("error 134: could not initialize db")
	}

	return &PostgresStore{db: db}, nil
}

func (pg *PostgresStore) CreateUser(u *UserType) error {
	var userid int //dont really need it
	row := pg.db.QueryRow("INSERT INTO users(name, email, pw, uid) VALUES ($1, $2, $3, $4) RETURNING id", u.name, u.email, u.pw, u.uid)
	err := row.Scan(&userid)

	if err != nil {

		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return fmt.Errorf("user %s already exists", u.email)
		}
		return errors.New("db error 430: could not add user")
	}
	return nil
}

func (*PostgresStore) UpdateUser(u *UserType) error {
	//TODO
	return nil
}

func (*PostgresStore) DeleteUser(u *UserType) error {
	//TODO
	return nil
}

func (pg *PostgresStore) ReadUserByEmail(req string) (*UserType, error) {
	row := pg.db.QueryRow("SELECT id, email, pw, name, uid from users where email=$1", req)

	var id int
	var email string
	var pw []byte
	var name string
	var uid string

	err := row.Scan(&id, &email, &pw, &name, &uid)
	if err != nil {
		return nil, fmt.Errorf("db error 431: could not read user %s", req)
	}

	return &UserType{
		id:    id,
		email: email,
		pw:    pw,
		name:  name,
		uid:   uid,
	}, nil
}
