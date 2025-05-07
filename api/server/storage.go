package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(u *UserType) error
	DeleteUser(u *UserType) error
	ReadUserByEmail(email string) (*UserType, error)
	UpdateUser(u *UserType) error
}

type TestStore struct {
	email string
	id    int
	pw    string
	name  string
}

func NewTestStore() (*TestStore, error) {
	return &TestStore{email: "nate@test.ch", id: 123, pw: "pw", name: "Nate"}, nil
}

func (*TestStore) CreateUser(u *UserType) error {
	fmt.Printf("Creating User %+v\n", u)
	return nil
}

func (*TestStore) UpdateUser(u *UserType) error {
	return nil
}

func (*TestStore) DeleteUser(u *UserType) error {
	return nil
}

func (s *TestStore) ReadUserByEmail(email string) (*UserType, error) {
	return &UserType{
		id:    s.id,
		email: s.email,
		pw:    s.pw,
		name:  s.name,
	}, nil
}

type PostgresStore struct {
	db *sql.DB
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err == nil {
		return
	}

	err = godotenv.Load("/run/secrets/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewPostGresStore() (*PostgresStore, error) {
	loadConfig()

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PW")

	// connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
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
	pw varchar not null,
	created_at timestamptz not null default now()
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
	fmt.Printf("Creating User %+v\n", u)

	var userid int //dont really need it
	row := pg.db.QueryRow("INSERT INTO users(name, email, pw) VALUES ($1, $2, $3) RETURNING id", u.name, u.email, u.pw)
	err := row.Scan(&userid)

	if err != nil {
		fmt.Println(err)
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
	row := pg.db.QueryRow("SELECT id, email, pw, name from users where email=$1", req)

	var id int
	var email string
	var pw string
	var name string
	err := row.Scan(&id, &email, &pw, &name)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("db error 431: could not read user")
	}

	return &UserType{
		id:    id,
		email: email,
		pw:    pw,
		name:  name,
	}, nil
}
