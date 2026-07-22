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

		assertResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setUpServeHttp: func(t *testing.T) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				request := httptest.NewRequest(http.MethodPost, "/v1/genpass", nil)

				app := api.NewEngine(nil)

				app.ServeHTTP(recorder, request)

				return recorder
			},
			assertResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.Contains(t, recorder.Body.String(), `{"password":`)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			recorder := testCase.setUpServeHttp(t)

			testCase.assertResponse(t, recorder)
		})
	}
}
