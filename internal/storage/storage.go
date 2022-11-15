// Package storage implements function to work with postgres sql.
package storage

import (
	//"context"
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log"
	"time"

	//"time"

	"github.com/google/uuid"
	// pq - driver for postgre
	_ "github.com/lib/pq"

	"github.com/MaximkaSha/gophkeeper/internal/models"
)

// Storage class.
type Storage struct {
	// DSN.
	ConnectionString string
	// DB handler.
	DB *sql.DB
}

// NewStorage Constrctor of storage. Needs dsn.
func NewStorage(dsn string) *Storage {
	s := new(Storage)
	s.ConnectionString = dsn
	err := s.initDB()
	if err != nil {
		log.Panic("database error!")
	}
	log.Println("DB Connected!")
	return s
}

// initDB - initialize database.
func (s *Storage) initDB() error {
	psqlconn := s.ConnectionString
	var err error
	s.DB, err = sql.Open("postgres", psqlconn)
	CheckError(err)
	err = s.DB.Ping()
	CheckError(err)
	/*err = s.CreateDBIfNotExist()
	CheckError(err) */
	err = s.CreateTableIfNotExist()
	CheckError(err)

	return err
}

// CheckError Checks and print databse error.
func CheckError(err error) {
	if err != nil {
		log.Fatalf("Database error: %s", err.Error())
	}
}

/*
func (s Storage) CreateDBIfNotExist() error {
	var query = `SELECT 'CREATE DATABASE gophkeeper'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname  = 'gophkeeper')`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := s.DB.ExecContext(ctx, query)
	return err
}
*/
// CreateTableIfNotExist Create DBs tables if needed
func (s Storage) CreateTableIfNotExist() error {
	var query = `
CREATE TABLE IF NOT EXISTS users
(
    email character varying(100) NOT NULL,
    password character varying(100) NOT NULL,
    id serial,
    secret bytea,
    CONSTRAINT users_pkey PRIMARY KEY (email)
);
CREATE TABLE IF NOT EXISTS ciphereddata
(
    data bytea,
    type character varying(100) NOT NULL,
    user_id serial,
    uuid uuid NOT NULL,
    CONSTRAINT ciphereddata_pkey PRIMARY KEY (uuid)
);
`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating  table", err.Error())
		return err
	}

	return err
}

// AddUser - insert user to database.
func (s Storage) AddUser(user models.User) error {
	var query = `
	INSERT INTO users (email,password,secret)
	VALUES ($1, $2, $3)`
	secret, err := s.GenSecretKey()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.DB.Exec(query, user.Email, user.Password, secret)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GenSecretKey Generate user privite key.
func (s Storage) GenSecretKey() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return ``, err
	}
	return string(data), nil
}

// GetUser - select user model from database.
func (s Storage) GetUser(user models.User) (models.User, error) {
	var query = `SELECT email,password,secret from users where email = $1`
	data := models.User{}
	err := s.DB.QueryRow(query, user.Email).Scan(&data.Email, &data.Password, &data.Secret)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}
	return data, nil
}

// AddCipheredData - insert ciphered data to database.
func (s Storage) AddCipheredData(data models.CipheredData) error {
	var query = `INSERT INTO ciphereddata (data, type, user_id, uuid)
		VALUES ($1, $2, (SELECT id from users where email = $3), $4)
		ON CONFLICT (uuid)
		DO UPDATE SET
		data = EXCLUDED.data,
		type = EXCLUDED.type`
	if data.ID == "" {
		data.ID = uuid.NewString()
	}
	_, err := s.DB.Exec(query, data.Data, data.Type, data.User, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetCipheredData - returns all users data from database by given user.
func (s Storage) GetCipheredData(email string) ([]models.CipheredData, error) {
	var query = `SELECT * from ciphereddata where user_id = (SELECT id from users where email = $1)`
	rows, err := s.DB.Query(query, email)
	if err != nil {
		log.Printf("Error %s when getting all data", err)
		return []models.CipheredData{}, err
	}
	defer rows.Close()
	data := []models.CipheredData{}
	counter := 0
	for rows.Next() {
		model := models.CipheredData{}
		if err := rows.Scan(&model.Data, &model.Type, &model.User, &model.ID); err != nil {
			log.Println(err)
			return []models.CipheredData{}, err
		}
		counter++
		data = append(data, model)
	}
	if counter == 0 {
		return []models.CipheredData{}, errors.New("no data for user")
	}
	return data, nil
}

// DelCiphereData - delete user data from database by given uuid.
func (s Storage) DelCiphereData(uuid string) error {
	var query = `DELETE from ciphereddata WHERE uuid = $1`
	_, err := s.DB.Exec(query, uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
