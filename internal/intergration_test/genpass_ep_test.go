package intergration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NghiaLeopard/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGenPassEP(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name           string
		setUpServeHttp func(t *testing.T) *httptest.ResponseRecorder

		ExpectedStatusCode   int
		ExpectedResponseBody string
	}{
		{
			name: "success",
			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodPost, "/genpass", nil)

				app := api.NewEngine()

				app.ServeHTTP(recorder, request)

				return recorder
			},
			ExpectedStatusCode:   http.StatusOK,
			ExpectedResponseBody: `{"password":}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			recorder := testCase.setUpServeHttp(t)

			assert.Equal(t, testCase.ExpectedStatusCode, recorder.Code)
			assert.Contains(t, testCase.ExpectedResponseBody, recorder.Body.String())
		})
	}
}
