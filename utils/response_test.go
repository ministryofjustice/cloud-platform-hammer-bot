package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	gin.SetMode("test")

	tests := []struct {
		name           string
		responseObj    Response
		expectedStatus int
		expectedObj    string
	}{
		{
			"GIVEN a param with the messsage key THEN return a response with the message",
			Response{200, []string{"test is good"}, nil, nil},
			200,
			"{\"message\":\"test is good\"}",
		},
		{
			"GIVEN a param with the data key THEN return a response with the data",
			Response{200, nil, nil, gin.H{"data": "test has data"}},
			200,
			"{\"data\":\"test has data\"}",
		},
		{
			"GIVEN a param with the error key THEN return a response with the error",
			Response{500, nil, []string{"test is error"}, nil},
			500,
			"{\"error\":\"test is error\"}",
		},
		{
			"GIVEN a param with no keys THEN return just the status",
			Response{200, nil, nil, nil},
			200,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			w := httptest.NewRecorder()
			c := gin.CreateTestContextOnly(w, r)
			SendResponse(c, tt.responseObj)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedObj, w.Body.String())
		})
	}
}
