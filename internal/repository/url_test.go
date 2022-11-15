package repository

import (
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/pkg/database"
)

var (
	accessID    = uuid.Must(uuid.NewUUID())
	accessToken = "qwerty"
)

type urlRepositoryTestSuite struct {
	testDB *database.TestDatabase
	suite.Suite
	urlRepository    Url
	accessRepository Access
}

func TestUrlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &urlRepositoryTestSuite{})
}

func (s *urlRepositoryTestSuite) SetupSuite() {
	var err error

	s.testDB, err = database.SetupTestDatabase()
	s.Require().NoError(err)

	if err := s.testDB.Migrate(true); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	s.urlRepository = NewUrl(s.testDB.DB)
	s.accessRepository = NewAccess(s.testDB.DB)

	access := &entity.Access{ID: accessID, Token: accessToken}

	err = s.accessRepository.Create(access)
	s.Require().NoError(err)
}

func (s *urlRepositoryTestSuite) SetupTest() {
	_, err := s.testDB.DB.Exec("TRUNCATE table urls CASCADE")
	s.Require().NoError(err)
}

func (s *urlRepositoryTestSuite) TestFindByCodeSuccess() {
	code := "qwerty"
	ID := uuid.Must(uuid.NewUUID())
	url := &entity.Url{ID: ID, Code: code, AccessID: accessID}

	err := s.urlRepository.Create(url)
	s.Require().NoError(err)

	actual, err := s.urlRepository.FindByCode(code)
	s.Require().NoError(err)

	s.Equal(ID, actual.ID)
	s.Equal(code, actual.Code)
}

func (s *urlRepositoryTestSuite) TestFindByCodeNotFoundNil() {
	actual, err := s.urlRepository.FindByCode("qwerty")
	s.Require().NoError(err)

	s.Nil(actual)
}

func (s *urlRepositoryTestSuite) TestRemoveOverdue() {
	code := "qwerty"
	url := &entity.Url{
		ID:         uuid.Must(uuid.NewUUID()),
		Code:       code,
		AliveUntil: time.Now().AddDate(0, 0, -1),
		AccessID:   accessID,
	}

	err := s.urlRepository.Create(url)
	s.Require().NoError(err)

	err = s.urlRepository.RemoveOverdue()
	s.Require().NoError(err)

	actual, err := s.urlRepository.FindByCode(code)
	s.NoError(err)
	s.Nil(actual)
}

func (s *urlRepositoryTestSuite) TestCreateSuccess() {
	code := "qwerty"

	url := &entity.Url{
		ID:       uuid.Must(uuid.NewUUID()),
		Code:     code,
		AccessID: accessID,
	}

	err := s.urlRepository.Create(url)
	s.Require().NoError(err)

	actual, err := s.urlRepository.FindByCode(code)
	s.NoError(err)
	s.Equal(url.ID, actual.ID)
	s.Equal(url.Code, actual.Code)
}

func (s *urlRepositoryTestSuite) TestCreateFailsOnDuplicateKey() {
	code := "qwerty"

	url := &entity.Url{
		ID:       uuid.Must(uuid.NewUUID()),
		Code:     code,
		AccessID: accessID,
	}

	err := s.urlRepository.Create(url)
	s.Require().NoError(err)

	err = s.urlRepository.Create(url)
	s.Error(err)
	s.Equal(err.Error(), "pq: duplicate key value violates unique constraint \"urls_pkey\"")
}

func (s *urlRepositoryTestSuite) TestDelete() {
	code := "qwerty"
	url := &entity.Url{
		ID:       uuid.Must(uuid.NewUUID()),
		Code:     code,
		AccessID: accessID,
	}

	err := s.urlRepository.Create(url)
	s.Require().NoError(err)

	err = s.urlRepository.Delete(url.ID)
	s.Require().NoError(err)

	actual, err := s.urlRepository.FindByCode(code)
	s.NoError(err)
	s.Nil(actual)
}

func (s *urlRepositoryTestSuite) TestFindByAccessID() {
	id1 := uuid.Must(uuid.NewUUID())
	id2 := uuid.Must(uuid.NewUUID())
	id3 := uuid.Must(uuid.NewUUID())
	code1 := "qwerty1"
	code2 := "qwerty2"
	code3 := "qwerty3"

	urls := []*entity.Url{
		{
			ID:       id1,
			Code:     code1,
			AccessID: accessID,
		},
		{
			ID:       id2,
			Code:     code2,
			AccessID: accessID,
		},
		{
			ID:       id3,
			Code:     code3,
			AccessID: accessID,
		},
	}

	for _, url := range urls {
		err := s.urlRepository.Create(url)
		s.Require().NoError(err)
	}

	actual, err := s.urlRepository.FindByAccessID(accessID, 0, 10)
	s.NoError(err)
	s.Equal(3, len(actual))
	s.Equal(accessID, actual[0].AccessID)
}

func (s *urlRepositoryTestSuite) TestFindByAccessIDLimit2() {
	id1 := uuid.Must(uuid.NewUUID())
	id2 := uuid.Must(uuid.NewUUID())
	id3 := uuid.Must(uuid.NewUUID())
	code1 := "qwerty1"
	code2 := "qwerty2"
	code3 := "qwerty3"

	urls := []*entity.Url{
		{
			ID:       id1,
			Code:     code1,
			AccessID: accessID,
		},
		{
			ID:       id2,
			Code:     code2,
			AccessID: accessID,
		},
		{
			ID:       id3,
			Code:     code3,
			AccessID: accessID,
		},
	}

	for _, url := range urls {
		err := s.urlRepository.Create(url)
		s.Require().NoError(err)
	}

	actual, err := s.urlRepository.FindByAccessID(accessID, 0, 2)
	s.NoError(err)
	s.Equal(2, len(actual))
	s.Equal(accessID, actual[0].AccessID)
}

func (s *urlRepositoryTestSuite) TestFindByAccessIDOffset2Limit2() {
	id1 := uuid.Must(uuid.NewUUID())
	id2 := uuid.Must(uuid.NewUUID())
	id3 := uuid.Must(uuid.NewUUID())
	code1 := "qwerty1"
	code2 := "qwerty2"
	code3 := "qwerty3"

	urls := []*entity.Url{
		{
			ID:       id1,
			Code:     code1,
			AccessID: accessID,
		},
		{
			ID:       id2,
			Code:     code2,
			AccessID: accessID,
		},
		{
			ID:       id3,
			Code:     code3,
			AccessID: accessID,
		},
	}

	for _, url := range urls {
		err := s.urlRepository.Create(url)
		s.Require().NoError(err)
	}

	actual, err := s.urlRepository.FindByAccessID(accessID, 2, 2)
	s.NoError(err)
	s.Equal(1, len(actual))
	s.Equal(accessID, actual[0].AccessID)
}
