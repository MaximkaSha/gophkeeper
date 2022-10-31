package storage

import (
	"context"
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

func NewStorage() *Storage {
	s := new(Storage)
	s.ConnectionString = "postgres://postgres:123456@127.0.0.1:5432/gophkeeper?sslmode=disable"
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
		log.Printf("Database error: %s", err)
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
	CREATE TABLE IF NOT EXISTS public.cc
(
    ccnum character varying(100) COLLATE pg_catalog."default",
    exp character varying(100) COLLATE pg_catalog."default",
    name character varying(100) COLLATE pg_catalog."default",
    cvv character varying(100) COLLATE pg_catalog."default",
    tag character varying COLLATE pg_catalog."default",
    uuid uuid NOT NULL,
    CONSTRAINT cc_pkey PRIMARY KEY (uuid)
);
CREATE TABLE IF NOT EXISTS public.data
(
    data bytea,
    tag character varying(100) COLLATE pg_catalog."default",
    uuid uuid NOT NULL,
    CONSTRAINT data_pkey PRIMARY KEY (uuid)
);
CREATE TABLE IF NOT EXISTS public.password
(
    login character varying(100) COLLATE pg_catalog."default",
    password character varying(100) COLLATE pg_catalog."default",
    tag character varying(100) COLLATE pg_catalog."default",
    uuid uuid NOT NULL,
    CONSTRAINT password_pkey PRIMARY KEY (uuid)
);
CREATE TABLE IF NOT EXISTS public.text
(
    text text COLLATE pg_catalog."default",
    tag character varying(100) COLLATE pg_catalog."default",
    uuid uuid NOT NULL,
    CONSTRAINT text_pkey PRIMARY KEY (uuid)
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

/*
type Storager interface {
	AddPassword(Password) error
	GetPassword(string) (Password, error)
	DelPassword(string) error
	UpdatePassword(string, Password) error
	GetAllPassword() ([]Password, error)

	AddData(Data) error
	GetData(string) (Data, error)
	DelData(string) error
	UpdateData(string, Data) error
	GetAllData() ([]Data, error)

	AddText(Text) error
	GetText(string) (Text, error)
	DelText(string) error
	UpdateText(string, Text) error
	GetAllText() ([]Text, error)

	AddCreditCard(CreditCard) error
	GetCreditCard(string) (CreditCard, error)
	DelCreditCard(string) error
	UpdateCreditCard(string, CreditCard)
	GetAllCreditCard() ([]CreditCard, error)
}
*/
func (s Storage) AddPassword(pass models.Password) error {
	var query = `
	INSERT INTO password (login,password,tag,uuid)
	VALUES ($1, $2, $3, $4)`
	_, err := s.DB.Exec(query, pass.Login, pass.Password, pass.Tag, uuid.New())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetPassword(uuid string) (models.Password, error) {
	var query = `SELECT login, password,tag,uuid from password WHERE uuid = $1`
	data := models.Password{}
	err := s.DB.QueryRow(query, uuid).Scan(&data.Login, &data.Password, &data.Tag, &data.ID)
	if err != nil {
		log.Println(err)
		return models.Password{}, err
	}
	return data, nil
}

func (s Storage) DelPassword(uuid string) error {
	var query = `DELETE from password WHERE uuid = $1`
	_, err := s.DB.Exec(query, uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) UpdatePassword(uuid string, pass models.Password) error {
	var query = `INSERT INTO password (login,password,tag,uuid) 
	VALUES ( $1, $2, $3, $4)
	ON CONFLICT (uuid)
	DO UPDATE SET
	login = EXCLUDED.login,
	password = EXCLUDED.password,
	tag = EXCLUDED.tag`
	_, err := s.DB.Exec(query, pass.Login, pass.Password, pass.Tag, pass.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetAllPassword() ([]models.Password, error) {
	var query = `SELECT * from password`
	rows, err := s.DB.Query(query)
	//	err = rows.Err()
	if err != nil {
		log.Printf("Error %s when getting all  data", err)
		return []models.Password{}, err
	}
	defer rows.Close()
	data := []models.Password{}
	counter := 0
	for rows.Next() {
		model := models.Password{}
		if err := rows.Scan(&model.Login, &model.Password, &model.Tag, &model.ID); err != nil {
			log.Println(err)
			return []models.Password{}, err
		}
		counter++
		data = append(data, model)
	}
	if counter == 0 {
		return []models.Password{}, errors.New("no data")
	}
	return data, nil
}
