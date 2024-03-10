package repository

import (
	"WatchHive/pkg/utils/models"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCheckUserAvailability(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		stub func(mock sqlmock.Sqlmock)
		want bool
	}{
		{
			name: "successful, user available",
			arg:  "arun7@gmail.com",
			stub: func(mock sqlmock.Sqlmock) {
				querry := "select count (*) from users where email='arun7@gmail.com'"
				mock.ExpectQuery(regexp.QuoteMeta(querry)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: true,
		}, {
			name: "failed, user not avilable",
			arg:  "arun1@gmail.com",
			stub: func(mock sqlmock.Sqlmock) {
				querry := "select count (*) from users where email ='arun1@gmail.com'"
				mock.ExpectQuery(regexp.QuoteMeta(querry)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		mockDb, mockSql, _ := sqlmock.New()
		DB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: mockDb,
		}), &gorm.Config{})
		userRepository := NewUserRepository(DB)
		tt.stub(mockSql)

		result := userRepository.CheckUserAvilability(tt.arg)
		assert.Equal(t, tt.want, result)
	}
}

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.UserDetails
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.UserDetailsResponse
		wantErr    error
	}{
		{
			name: "Successfully user signed up",
			args: args{
				input: models.UserDetails{
					Name:     "Rahul",
					Email:    "rahulchacko888@gmail.com",
					Password: "12345",
					Phone:    "9867327710",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users \(name, email, password, phone\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, name, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("Rahul", "rahulchacko888@gmail.com", "12345", "9867327710").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
						AddRow(1, "Rahul", "rahulchacko888@gmail.com", "9867327710"))
			},
			want: models.UserDetailsResponse{
				Id:    1,
				Name:  "Rahul",
				Email: "rahulchacko888@gmail.com",
				Phone: "9867327710",
			},
			wantErr: nil,
		},
		{
			name: "Error signing up user",
			args: args{
				input: models.UserDetails{},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users \(name, email, password, phone\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, name, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("", "", "", "").
					WillReturnError(errors.New("email should be unique"))
			},
			want:    models.UserDetailsResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSql)
			u := NewUserRepository(gormDB)
			got, err := u.UserSignup(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})

	}
}
