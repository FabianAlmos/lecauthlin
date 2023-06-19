package util

import (
	"auth/response"
	"auth/util"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertTokenValidationResult(testCase TestCaseValidate, t *testing.T, gotError error, gotClaims *util.JwtCustomClaims) {
	t.Helper()

	if testCase.WantError {
		assert.Error(t, gotError)
		assert.Contains(t, gotError.Error(), testCase.WantErrorMsg)
	} else {
		assert.NoError(t, gotError)
		assert.Equal(t, testCase.WantId, gotClaims.ID)
	}
}

func AssertUserProfileResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var r response.UserResponse
	err := json.Unmarshal([]byte(recorder.Body.String()), &r)

	if assert.NoError(t, err) {
		assert.Equal(t, response.UserResponse{
			ID:    1,
			Name:  "Test 1",
			Email: "test-1@example.com",
		}, r)
	}
}
