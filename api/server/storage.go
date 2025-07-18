package server

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(u *UserType) error
	DeleteUser(u *UserType) error
	ReadUserByEmail(email string) (*UserType, error)
	ReadUserByUid(uid string) (*UserType, error)
	UpdateUser(u *UserType) error

	CreateRefreshToken(t *RefreshTokenType) error
	DeleteRefreshToken(t *RefreshTokenType) error
	ReadRefreshTokenByToken(token string) (*RefreshTokenType, error)
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
		return nil, err
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
		return nil, errors.New("error 133a: could not initialize db")
	}

	_, err = db.Query(`
	CREATE TABLE 
	IF NOT EXISTS 
	refresh_tokens(
	id SERIAL PRIMARY KEY,
	user_uid varchar not null,
	token varchar not null,
	expires timestamptz not null,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
	)
	`)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error 133b: could not initialize db")
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

func (pg *PostgresStore) ReadUserByUid(uid string) (*UserType, error) {
	row := pg.db.QueryRow("SELECT id, email, pw, name, uid from users where uid=$1", uid)

	var id int
	var email string
	var pw []byte
	var name string
	var uidRes string

	err := row.Scan(&id, &email, &pw, &name, &uidRes)
	if err != nil {
		return nil, fmt.Errorf("db error 441: could not read user %s", uid)
	}

	return &UserType{
		id:    id,
		email: email,
		pw:    pw,
		name:  name,
		uid:   uidRes,
	}, nil
}

func (pg *PostgresStore) CreateRefreshToken(t *RefreshTokenType) error {
	var id int //dont really need it
	row := pg.db.QueryRow("INSERT INTO refresh_tokens(user_uid, token, expires) VALUES ($1, $2, $3) RETURNING id", t.userUid, t.token, t.expiration)
	err := row.Scan(&id)

	if err != nil {
		return err
		// return errors.New("db error 610: could not add new refresh_token")
	}
	return nil
}

func (pg *PostgresStore) DeleteRefreshToken(t *RefreshTokenType) error {
	_, err := pg.db.Query("DELETE FROM refresh_token rt where rt.id = $1", t.id)
	if err != nil {
		return fmt.Errorf("db error 620: could not delete refresh_token %v", t.id)
	}
	return nil
}

func (pg *PostgresStore) ReadRefreshTokenByToken(token string) (*RefreshTokenType, error) {
	row := pg.db.QueryRow("SELECT id, user_uid, token, expires from refresh_tokens where token = $1", token)

	var id int
	var tokenRes string
	var userUid string
	var expiration time.Time

	err := row.Scan(&id, &userUid, &tokenRes, &expiration)
	if err != nil {
		return nil, fmt.Errorf("db error 630: could not read refresh_token %s", token)
	}

	return &RefreshTokenType{
		id:         id,
		token:      tokenRes,
		userUid:    userUid,
		expiration: expiration,
	}, nil
}
