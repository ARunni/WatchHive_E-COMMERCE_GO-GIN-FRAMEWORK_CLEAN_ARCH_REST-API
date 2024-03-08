package repository

import (
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
