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
CREATE TABLE IF NOT EXISTS public.users
(
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    password character varying(100) COLLATE pg_catalog."default" NOT NULL,
    secret character varying(100) COLLATE pg_catalog."default",
	id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    CONSTRAINT users_pkey PRIMARY KEY (email)
);
CREATE TABLE IF NOT EXISTS public.ciphereddata
(
    data bytea,
    "type " character varying(100) COLLATE pg_catalog."default" NOT NULL,
    user_id bigint NOT NULL DEFAULT nextval('ciphereddata_user_id_seq'::regclass),
    uuid uuid
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

// Password storage endpoints.
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

// Data storage block.
func (s Storage) AddData(data models.Data) error {
	var query = `
	INSERT INTO data (data,tag,uuid)
	VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, data.Data, data.Tag, uuid.New())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetData(uuid string) (models.Data, error) {
	var query = `SELECT data,tag,uuid from data WHERE uuid = $1`
	data := models.Data{}
	err := s.DB.QueryRow(query, uuid).Scan(&data.Data, &data.Tag, &data.ID)
	if err != nil {
		log.Println(err)
		return models.Data{}, err
	}
	return data, nil
}

func (s Storage) DelData(uuid string) error {
	var query = `DELETE from data WHERE uuid = $1`
	_, err := s.DB.Exec(query, uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) UpdateData(uuid string, data models.Data) error {
	var query = `INSERT INTO data (data,tag,uuid) 
	VALUES ( $1, $2, $3)
	ON CONFLICT (uuid)
	DO UPDATE SET
	data = EXCLUDED.data,
	tag = EXCLUDED.tag`
	_, err := s.DB.Exec(query, data.Data, data.Tag, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetAllData() ([]models.Data, error) {
	var query = `SELECT * from data`
	rows, err := s.DB.Query(query)
	//	err = rows.Err()
	if err != nil {
		log.Printf("Error %s when getting all  data", err)
		return []models.Data{}, err
	}
	defer rows.Close()
	data := []models.Data{}
	counter := 0
	for rows.Next() {
		model := models.Data{}
		if err := rows.Scan(&model.Data, &model.Tag, &model.ID); err != nil {
			log.Println(err)
			return []models.Data{}, err
		}
		counter++
		data = append(data, model)
	}
	if counter == 0 {
		return []models.Data{}, errors.New("no data")
	}
	return data, nil
}

// Text storage block.
func (s Storage) AddText(data models.Text) error {
	var query = `
	INSERT INTO text (text,tag,uuid)
	VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, data.Data, data.Tag, uuid.New())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetText(uuid string) (models.Text, error) {
	var query = `SELECT text,tag,uuid from text WHERE uuid = $1`
	data := models.Text{}
	err := s.DB.QueryRow(query, uuid).Scan(&data.Data, &data.Tag, &data.ID)
	if err != nil {
		log.Println(err)
		return models.Text{}, err
	}
	return data, nil
}

func (s Storage) DelText(uuid string) error {
	var query = `DELETE from text WHERE uuid = $1`
	_, err := s.DB.Exec(query, uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) UpdateText(uuid string, data models.Text) error {
	var query = `INSERT INTO text (text,tag,uuid) 
	VALUES ( $1, $2, $3)
	ON CONFLICT (uuid)
	DO UPDATE SET
	text = EXCLUDED.text,
	tag = EXCLUDED.tag`
	_, err := s.DB.Exec(query, data.Data, data.Tag, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetAllText() ([]models.Text, error) {
	var query = `SELECT * from text`
	rows, err := s.DB.Query(query)
	//	err = rows.Err()
	if err != nil {
		log.Printf("Error %s when getting all  text", err)
		return []models.Text{}, err
	}
	defer rows.Close()
	data := []models.Text{}
	counter := 0
	for rows.Next() {
		model := models.Text{}
		if err := rows.Scan(&model.Data, &model.Tag, &model.ID); err != nil {
			log.Println(err)
			return []models.Text{}, err
		}
		counter++
		data = append(data, model)
	}
	if counter == 0 {
		return []models.Text{}, errors.New("no text")
	}
	return data, nil
}

// CC storage block.
func (s Storage) AddCreditCard(data models.CreditCard) error {
	var query = `
	INSERT INTO cc (ccnum,exp,name,cvv,tag,uuid)
	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.DB.Exec(query, data.CardNum, data.Exp, data.Name, data.CVV, data.Tag, uuid.New())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetCreditCard(uuid string) (models.CreditCard, error) {
	var query = `SELECT ccnum,exp,name,cvv,tag,uuid from cc WHERE uuid = $1`
	data := models.CreditCard{}
	err := s.DB.QueryRow(query, uuid).Scan(&data.CardNum, &data.Exp, &data.Name, &data.CVV, &data.Tag, &data.ID)
	if err != nil {
		log.Println(err)
		return models.CreditCard{}, err
	}
	return data, nil
}

func (s Storage) DelCreditCard(uuid string) error {
	var query = `DELETE from cc WHERE uuid = $1`
	_, err := s.DB.Exec(query, uuid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) UpdateCreditCard(uuid string, data models.CreditCard) error {
	var query = `INSERT INTO cc (ccnum,exp,name,cvv,tag,uuid) 
	VALUES ( $1, $2, $3, $4, $5, $6)
	ON CONFLICT (uuid)
	DO UPDATE SET
	ccnum = EXCLUDED.ccnum,
	exp = EXCLUDED.exp,
	name = EXCLUDED.name,
	cvv = EXCLUDED.cvv,
	tag = EXCLUDED.tag`
	_, err := s.DB.Exec(query, data.CardNum, data.Exp, data.Name, data.CVV, data.Tag, data.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s Storage) GetAllCreditCard() ([]models.CreditCard, error) {
	var query = `SELECT * from cc`
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Printf("Error %s when getting all cc", err)
		return []models.CreditCard{}, err
	}
	defer rows.Close()
	data := []models.CreditCard{}
	counter := 0
	for rows.Next() {
		model := models.CreditCard{}
		if err := rows.Scan(&model.CardNum, &model.Exp, &model.Name, &model.CVV, &model.Tag, &model.ID); err != nil {
			log.Println(err)
			return []models.CreditCard{}, err
		}
		counter++
		data = append(data, model)
	}
	if counter == 0 {
		return []models.CreditCard{}, errors.New("no cc")
	}
	return data, nil
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
