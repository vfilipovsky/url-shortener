package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	mockAction "github.com/vfilipovsky/url-shortener/internal/action/url/mock"
	"github.com/vfilipovsky/url-shortener/internal/entity"
	"github.com/vfilipovsky/url-shortener/pkg/config"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type urlHandlerE2eTestSuite struct {
	suite.Suite
	handler               *urlHandler
	getUrlsByAccessAction *mockAction.MockGetAllByAccess
	createUrlAction       *mockAction.MockCreate
	deleteUrlAction       *mockAction.MockDelete
	getUrlAction          *mockAction.MockGet
	infoUrlAction         *mockAction.MockInfo
}

func TestUrlE2eTestSuite(t *testing.T) {
	suite.Run(t, &urlHandlerE2eTestSuite{})
}

func (s *urlHandlerE2eTestSuite) SetupSuite() {
	err := godotenv.Load("../../.env.test")
	s.Require().NoError(err)
	cfg := config.Init()
	logger.NewInstance(cfg.Logger)

	ctrl := gomock.NewController(s.Suite.T())
	s.getUrlsByAccessAction = mockAction.NewMockGetAllByAccess(ctrl)
	s.createUrlAction = mockAction.NewMockCreate(ctrl)
	s.deleteUrlAction = mockAction.NewMockDelete(ctrl)
	s.getUrlAction = mockAction.NewMockGet(ctrl)
	s.infoUrlAction = mockAction.NewMockInfo(ctrl)

	s.handler = newUrl(
		cfg.Http,
		s.createUrlAction,
		s.deleteUrlAction,
		s.getUrlAction,
		s.infoUrlAction,
		s.getUrlsByAccessAction,
	)
}

func (s *urlHandlerE2eTestSuite) TestGetUrlsByAccess() {
	token := "qwerty1234"
	reader := strings.NewReader(fmt.Sprintf(`{"access_token":"%s"}`, token))
	r := httptest.NewRequest(http.MethodPost, "/api/v1/url?offset=0&limit=100", reader)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	urls := []*entity.Url{
		{
			Code: "qwerty",
			Url:  "https://google.com",
		},
		{
			Code: "zxcvb",
			Url:  "https://bing.com",
		},
	}

	s.getUrlsByAccessAction.EXPECT().Run(token, 0, 100).Return(urls, nil)

	s.handler.GetUrlsByAccess(w, r)

	var actualUrls []*entity.Url
	err := json.NewDecoder(w.Result().Body).Decode(&actualUrls)
	s.NoError(err)

	s.Equal(urls, actualUrls)
	s.Equal(http.StatusOK, w.Result().StatusCode)
}

func (s *urlHandlerE2eTestSuite) TestCreate() {
	token := "qwerty1234"
	isSecured := false
	urlString := "https://google.com"
	aliveUntilString := "2021-12-23T23:23:23Z"
	aliveUntil, err := time.Parse(time.RFC3339, "2021-12-23T23:23:23Z")

	createJson := fmt.Sprintf(`{
		"access_token": "%s",
		"alive_until": "%s",
		"is_secured": %s,
		"url": "%s"
	  }`, token, aliveUntilString, strconv.FormatBool(isSecured), urlString)

	s.Require().NoError(err)

	reader := strings.NewReader(createJson)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/url", reader)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	url := &entity.Url{
		IsSecured:  isSecured,
		Url:        urlString,
		AliveUntil: aliveUntil,
	}

	s.createUrlAction.EXPECT().Run(token, urlString, isSecured, aliveUntil).Return(url, nil)

	s.handler.Create(w, r)

	var actualUrl *entity.Url
	err = json.NewDecoder(w.Result().Body).Decode(&actualUrl)
	s.NoError(err)

	s.Equal(url.Url, actualUrl.Url)
	s.Equal(url.IsSecured, actualUrl.IsSecured)
	s.Equal(http.StatusOK, w.Result().StatusCode)
}
