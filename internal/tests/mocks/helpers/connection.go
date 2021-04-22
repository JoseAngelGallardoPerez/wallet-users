package helpers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

type dbMock struct {
	mock     sqlmock.Sqlmock
	gormMock *gorm.DB
}

var DbMock *dbMock

func init() {
	DbMock = new(dbMock)
}

func (s *dbMock) GetDbMock() sqlmock.Sqlmock {
	if s.mock == nil {
		db, mock, _ := sqlmock.New()
		gormDb, _ := gorm.Open("mysql", db)

		s.mock = mock
		s.gormMock = gormDb
	}

	return s.mock
}

func (s *dbMock) GetGormMock() *gorm.DB {
	if s.gormMock == nil {
		_ = s.GetDbMock()
	}

	return s.gormMock
}
