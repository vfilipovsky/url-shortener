package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/internal/repository"
	"github.com/vfilipovsky/url-shortener/internal/server"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/database"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type urlE2eTestSuite struct {
	suite.Suite
	testDB           *database.TestDatabase
	urlRepository    repository.Url
	accessRepository repository.Access
	serverUrl        string
}

func TestUrlE2eTestSuite(t *testing.T) {
	suite.Run(t, &urlE2eTestSuite{})
}

func (s *urlE2eTestSuite) SetupSuite() {
	var err error

	err = godotenv.Load("../.env.test")
	s.Require().NoError(err)
	cfg := config.Init()
	logger.NewInstance(cfg.Logger)

	s.testDB, err = database.SetupTestDatabase()
	s.Require().NoError(err)

	if err := s.testDB.Migrate(true); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	s.urlRepository = repository.NewUrl(s.testDB.DB)
	s.accessRepository = repository.NewAccess(s.testDB.DB)

	s.serverUrl = fmt.Sprintf("http://%s:%s", cfg.Http.Host, cfg.Http.Port)

	server := server.New(cfg, s.testDB.DB)
	go server.Run()
}

func (s *urlE2eTestSuite) SetupTest() {
	_, err := s.testDB.DB.Exec("TRUNCATE table accesses CASCADE")
	s.Require().NoError(err)
}

func (s *urlE2eTestSuite) TestGetUrlsByAccess() {
	id1 := uuid.Must(uuid.NewUUID())
	id2 := uuid.Must(uuid.NewUUID())
	id3 := uuid.Must(uuid.NewUUID())
	accessID1 := uuid.Must(uuid.NewUUID())
	accessID2 := uuid.Must(uuid.NewUUID())

	accesses := []*entity.Access{
		{
			ID:       accessID1,
			Token:    "qwerty1234",
			IsActive: true,
		},
		{
			ID:       accessID2,
			Token:    "1234qwerty",
			IsActive: true,
		},
	}

	urls := []*entity.Url{
		{
			ID:        id1,
			Code:      "asdf",
			AccessID:  accessID1,
			IsSecured: false,
		},
		{
			ID:        id2,
			Code:      "fdsa",
			AccessID:  accessID1,
			IsSecured: true,
			Pin:       "1234",
		},
		{
			ID:        id3,
			Code:      "zxcv",
			AccessID:  accessID2,
			IsSecured: true,
			Pin:       "4321",
		},
	}

	for _, access := range accesses {
		err := s.accessRepository.Create(access)
		s.Require().NoError(err)
	}

	for _, url := range urls {
		err := s.urlRepository.Create(url)
		s.Require().NoError(err)
	}

	reader := strings.NewReader(`{"access_token":"qwerty1234"}`)
	r, err := http.NewRequest(http.MethodPost, s.serverUrl+"/api/v1/url/access/list", reader)
	s.Require().NoError(err)

	r.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(r)
	s.Require().NoError(err)
	defer res.Body.Close()

	var resUrls []*entity.Url

	byteRes, err := io.ReadAll(res.Body)
	s.Require().NoError(err)

	err = json.Unmarshal(byteRes, &resUrls)
	s.Require().NoError(err)

	s.Equal(2, len(resUrls))
	s.Equal(http.StatusOK, res.StatusCode)
}
