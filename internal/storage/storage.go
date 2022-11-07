package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/MaximkaSha/gophkeeper/internal/models"
)

type Storage struct {
	ConnectionString string
	DB               *sql.DB
}

func NewStorage(dsn string) *Storage {
	s := new(Storage)
	s.ConnectionString = dsn
	err := s.InitDB()
	if err != nil {
		log.Panic("database error!")
	}
	log.Println("DB Connected!")
	return s
}

func (s *Storage) InitDB() error {
	psqlconn := s.ConnectionString
	var err error
	s.DB, err = sql.Open("postgres", psqlconn)
	CheckError(err)
	err = s.DB.Ping()
	CheckError(err)
	err = s.CreateDBIfNotExist()
	CheckError(err)
	err = s.CreateTableIfNotExist()
	CheckError(err)

	return err
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Database error: %s", err.Error())
	}
}

func (s Storage) CreateDBIfNotExist() error {
	var query = `SELECT 'CREATE DATABASE gophkeeper'
	WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname  = 'gophkeeper')`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := s.DB.ExecContext(ctx, query)
	return err
}

func (s Storage) CreateTableIfNotExist() error {
	var query = `
CREATE TABLE IF NOT EXISTS public.users
(
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    password character varying(100) COLLATE pg_catalog."default" NOT NULL,
    id serial,
    secret bytea,
    CONSTRAINT users_pkey PRIMARY KEY (email)
);
CREATE TABLE IF NOT EXISTS public.ciphereddata
(
    data bytea,
    type character varying(100) COLLATE pg_catalog."default" NOT NULL,
    user_id serial,
    uuid uuid NOT NULL,
    CONSTRAINT ciphereddata_pkey PRIMARY KEY (uuid)
);
`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating  table", err)
		return err
	}

	return err
}

// User group.
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

func (s Storage) GenSecretKey() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return ``, err
	}
	return string(data), nil
}

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

func (s Storage) AddCipheredData(data models.CipheredData) error {
	var query = `INSERT INTO ciphereddata (data, type, user_id, uuid)
		VALUES ($1, $2, (SELECT id from users where email = $3), $4)
		ON CONFLICT (uuid)
		DO UPDATE SET
		data = EXCLUDED.data,
		type = EXCLUDED.type`
	if data.Id == "" {
		data.Id = uuid.NewString()
	}
	_, err := s.DB.Exec(query, data.Data, data.Type, data.User, data.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetCipheredData(in models.CipheredData) ([]models.CipheredData, error) {
	var query = `SELECT * from ciphereddata where user_id = (SELECT id from users where email = $1)`
	rows, err := s.DB.Query(query, in.User)
	if err != nil {
		log.Printf("Error %s when getting all cc", err)
		return []models.CipheredData{}, err
	}
	defer rows.Close()
	data := []models.CipheredData{}
	counter := 0
	for rows.Next() {
		model := models.CipheredData{}
		if err := rows.Scan(&model.Data, &model.Type, &model.User, &model.Id); err != nil {
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

func (s Storage) DelCiphereData(in models.CipheredData) error {
	var query = `DELETE from ciphereddata WHERE uuid = $1`
	_, err := s.DB.Exec(query, in.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
