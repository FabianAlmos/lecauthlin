package handler

import (
	"auth/config"
	"auth/repository"
	"auth/service"
	"auth/test/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const userID = 1

type UserHandlerTestSuite struct {
	suite.Suite
	accessToken string
	userHandler *UserHandler
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	cfg := &config.Config{
		AccessSecret:          "secret",
		AccessLifetimeMinutes: 1,
	}
	tokenService := service.NewTokenService(cfg)

	suite.accessToken, _ = tokenService.GenerateAccessToken(userID)
	suite.userHandler = NewUserHandler(repository.NewUserRepositoryMock(), tokenService)
}

func (suite *UserHandlerTestSuite) TearDownSuite() {}

func (suite *UserHandlerTestSuite) SetupTest() {}

func (suite *UserHandlerTestSuite) TearDownTest() {}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) TestWalkUserHandlerGetProfile() {
	t := suite.T()
	handlerFunc := suite.userHandler.GetProfile

	cases := []util.TestCaseHandler{
		{
			TestName: "Successfully get user profile",
			Request: util.Request{
				Method:    http.MethodGet,
				Url:       "/profile",
				AuthToken: suite.accessToken,
			},
			HandlerFunc: handlerFunc,
			Want: util.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "test-1@example.com",
			},
		},
		{
			TestName: "Unauthorized getting user profile",
			Request: util.Request{
				Method:    http.MethodGet,
				Url:       "/profile",
				AuthToken: "",
			},
			HandlerFunc: handlerFunc,
			Want: util.ExpectedResponse{
				StatusCode: 401,
				BodyPart:   "Invalid credentials",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			request, recorder := util.PrepareHandlerTestCase(test)

			test.HandlerFunc(recorder, request)

			assert.Contains(t, recorder.Body.String(), test.Want.BodyPart)

			if assert.Equal(t, recorder.Code, test.Want.StatusCode) {
				if recorder.Code == http.StatusOK {
					util.AssertUserProfileResponse(t, recorder)
				}
			}
		})
	}
}
