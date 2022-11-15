package repository

import (
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/pkg/database"
)

type accessRepositoryTestSuite struct {
	testDB *database.TestDatabase
	suite.Suite
	accessRepository Access
}

func TestAccessRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &accessRepositoryTestSuite{})
}

func (s *accessRepositoryTestSuite) SetupSuite() {
	var err error

	s.testDB, err = database.SetupTestDatabase()
	s.Require().NoError(err)

	if err := s.testDB.Migrate(true); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	s.accessRepository = NewAccess(s.testDB.DB)
}

func (s *accessRepositoryTestSuite) SetupTest() {
	_, err := s.testDB.DB.Exec("TRUNCATE table accesses CASCADE")
	s.Require().NoError(err)
}

func (s *accessRepositoryTestSuite) TestCreateSuccess() {
	id := uuid.Must(uuid.NewUUID())
	token := "qwerty1234"
	access := &entity.Access{ID: id, Token: token}

	err := s.accessRepository.Create(access)
	s.Require().NoError(err)

	a, err := s.accessRepository.FindByToken(token)
	s.NoError(err)
	s.Equal(token, a.Token)
}

func (s *accessRepositoryTestSuite) TestCreateFailsOnDuplicateKey() {
	id := uuid.Must(uuid.NewUUID())
	token := "qwerty1234"
	access := &entity.Access{ID: id, Token: token}

	err := s.accessRepository.Create(access)
	s.Require().NoError(err)

	err = s.accessRepository.Create(access)
	s.Error(err)
	s.Equal("pq: duplicate key value violates unique constraint \"accesses_pkey\"", err.Error())
}

func (s *accessRepositoryTestSuite) TestFindByTokenSuccess() {
	id := uuid.Must(uuid.NewUUID())
	token := "qwerty1234"
	access := &entity.Access{ID: id, Token: token}

	err := s.accessRepository.Create(access)
	s.Require().NoError(err)

	actual, err := s.accessRepository.FindByToken(token)
	s.NoError(err)
	s.NotNil(actual)
	s.Equal(token, actual.Token)
	s.Equal(id, actual.ID)
}

func (s *accessRepositoryTestSuite) TestFindByTokenNotFoundNil() {
	actual, err := s.accessRepository.FindByToken("qwerty")
	s.NoError(err)
	s.Nil(actual)
}
