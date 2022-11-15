package storage

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaximkaSha/gophkeeper/internal/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddUser(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name string
		s    Storage
		args args
	}{
		{
			name: "pos 1",
			s:    Storage{},
			args: args{
				user: models.User{
					Email:    "test@test.com",
					Password: "pass",
					Secret:   []byte(""),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tt.s.DB = db
			mock.ExpectExec("INSERT INTO").WithArgs(tt.args.user.Email, tt.args.user.Password, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			err = tt.s.AddUser(tt.args.user)
			require.NoError(t, err)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			mock.ExpectExec("INSERT INTO").WithArgs(tt.args.user.Email, tt.args.user.Password, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(errors.New("no"))
			err = tt.s.AddUser(tt.args.user)
			require.Error(t, err)
		})
	}
}

func TestStorage_GenSecretKey(t *testing.T) {
	tests := []struct {
		name    string
		s       Storage
		wantErr bool
	}{
		{
			name:    "pos 1",
			s:       Storage{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GenSecretKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GenSecretKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got2, err := tt.s.GenSecretKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GenSecretKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == got2 {
				t.Errorf("Storage.GenSecretKey() equal keys")
			}
		})
	}
}

func TestStorage_GetUser(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name string
		s    Storage
		args args
	}{
		{
			name: "test 1",
			s:    Storage{},
			args: args{
				user: models.User{
					Email:    "test@test.com",
					Password: "pass",
					Secret:   []byte("secret"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tt.s.DB = db
			mockUserRows := sqlmock.NewRows([]string{"email", "password", "secret"}).AddRow(
				"test@test.com", "passNew", "newSecret",
			)
			mock.ExpectQuery("SELECT email,password,secret from users where email = ?").WithArgs(tt.args.user.Email).WillReturnRows(mockUserRows)
			_, err = tt.s.GetUser(tt.args.user)
			require.NoError(t, err)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			mock.ExpectQuery("SELECT email,password,secret from users where email = ?").WithArgs("no data").WillReturnError(errors.New("no data"))
			tt.args.user.Email = "no data"
			_, err = tt.s.GetUser(tt.args.user)
			require.Error(t, err)
		})
	}
}

func TestStorage_AddCipheredData(t *testing.T) {
	type args struct {
		data models.CipheredData
	}
	tests := []struct {
		name    string
		s       Storage
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			s:    Storage{},
			args: args{
				data: models.CipheredData{
					Type: "CC",
					Data: []byte("123"),
					User: "test@test.com",
					ID:   "1111-1111-111-11",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tt.s.DB = db

			mock.ExpectExec("INSERT INTO").WithArgs(tt.args.data.Data, tt.args.data.Type, tt.args.data.User, tt.args.data.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			err = tt.s.AddCipheredData(tt.args.data)
			require.NoError(t, err)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			mock.ExpectExec("INSERT INTO").WithArgs(tt.args.data.Data, "no data", tt.args.data.User, tt.args.data.ID).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("no data"))
			tt.args.data.Type = "no data"
			err = tt.s.AddCipheredData(tt.args.data)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			require.Error(t, err)

		})
	}
}

func TestStorage_GetCipheredData(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		s       Storage
		args    args
		want    []models.CipheredData
		wantErr bool
	}{
		{
			name: "test 1",
			s:    Storage{},
			args: args{
				email: "test@test.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tt.s.DB = db
			mockDataRows := sqlmock.NewRows([]string{"email", "password", "secret", "data"}).AddRow(
				"test@test.com", "passNew", "newSecret", "data",
			)
			mock.ExpectQuery("SELECT (.+) from ciphereddata where user_id").WithArgs(tt.args.email).WillReturnRows(mockDataRows)
			_, err = tt.s.GetCipheredData(tt.args.email)
			require.NoError(t, err)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			mock.ExpectQuery("SELECT (.+) from ciphereddata where user_id").WithArgs(tt.args.email).WillReturnRows(mockDataRows).WillReturnError(errors.New("no data"))
			_, err = tt.s.GetCipheredData(tt.args.email)
			require.Error(t, err)
			mock.ExpectQuery("SELECT (.+) from ciphereddata where user_id").WithArgs(tt.args.email).WillReturnRows()
		})
	}
}

func TestStorage_DelCiphereData(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		s       Storage
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			s:    Storage{},
			args: args{
				uuid: "111-11-11-111",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tt.s.DB = db

			mock.ExpectExec("DELETE from ciphereddata WHERE").WithArgs(tt.args.uuid).WillReturnResult(sqlmock.NewResult(1, 1))
			err = tt.s.DelCiphereData(tt.args.uuid)
			require.NoError(t, err)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			mock.ExpectExec("DELETE from ciphereddata WHERE").WithArgs(tt.args.uuid).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(errors.New("no data"))
			err = tt.s.DelCiphereData(tt.args.uuid)
			require.Error(t, err)
		})
	}
}

func TestStorage_initDB(t *testing.T) {
	tests := []struct {
		name    string
		s       *Storage
		wantErr bool
	}{
		{
			name:    "test",
			s:       &Storage{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.initDB(); (err != nil) != tt.wantErr {
				t.Errorf("Storage.initDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
