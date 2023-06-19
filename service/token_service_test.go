package service

import (
	"auth/config"
	"auth/test/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const userID = 1

type TokenServiceTestSuite struct {
	suite.Suite
	cfg          *config.Config
	tokenService *TokenService
}

func TestTokenServiceSuite(t *testing.T) {
	suite.Run(t, new(TokenServiceTestSuite))
}

func (suite *TokenServiceTestSuite) SetupSuite() {
	suite.cfg = &config.Config{
		AccessSecret:           "a_secret",
		RefreshSecret:          "r_secret",
		AccessLifetimeMinutes:  1,
		RefreshLifetimeMinutes: 2,
	}
	suite.tokenService = NewTokenService(suite.cfg)
}

func (suite *TokenServiceTestSuite) TearDownSuite() {

}

func (suite *TokenServiceTestSuite) SetupTest() {

}

func (suite *TokenServiceTestSuite) TearDownTest() {

}

func (suite *TokenServiceTestSuite) TestGetTokenFromBearerString() {
	testCases := []util.TestCaseBearerToken{
		{
			BearerString: "Bearer test_token",
			Want:         "test_token",
		},
		{
			BearerString: "Bearer",
			Want:         "",
		},
		{
			BearerString: "Beare test_token",
			Want:         "",
		},
	}

	for _, testCase := range testCases {
		suite.T().Run("", func(t *testing.T) {
			got := suite.tokenService.GetTokenFromBearerString(testCase.BearerString)
			assert.Equal(t, testCase.Want, got)
		})
	}
}

func (suite *TokenServiceTestSuite) TestValidateAccessToken() {
	tokenString, _ := suite.tokenService.GenerateAccessToken(userID)
	refreshTokenString, _ := suite.tokenService.GenerateRefreshToken(userID)
	invalidTokenString := tokenString + "f"

	suite.cfg.AccessLifetimeMinutes = 0
	expiredTokenString, _ := suite.tokenService.GenerateAccessToken(userID)

	testCases := []util.TestCaseValidate{
		{
			Name:         "Valid token validated successfully",
			Token:        tokenString,
			WantError:    false,
			WantErrorMsg: "",
			WantId:       userID,
		},
		{
			Name:         "Valid refresh token is not accepted",
			Token:        refreshTokenString,
			WantError:    true,
			WantErrorMsg: "signature is invalid",
			WantId:       0,
		},
		{
			Name:         "Invalid token is not accepted",
			Token:        invalidTokenString,
			WantError:    true,
			WantErrorMsg: "signature is invalid",
			WantId:       0,
		},
		{
			Name:         "Expired token is not accepted",
			Token:        expiredTokenString,
			WantError:    true,
			WantErrorMsg: "token is expired",
			WantId:       0,
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.Name, func(t *testing.T) {
			gotClaims, gotError := suite.tokenService.ValidateAccessToken(testCase.Token)

			util.AssertTokenValidationResult(testCase, t, gotError, gotClaims)
		})
	}
}
