package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		expectedLength int
		expectedError  error
	}{
		{
			name:           "Success",
			expectedLength: 12,
			expectedError:  nil,
		},
		{
			name:           "Success with custom length",
			expectedLength: 16,
			expectedError:  nil,
		},
		{
			name:           "Success with custom length less than 8",
			expectedLength: 100000,
			expectedError:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			t.Parallel()

			testService := NewGenPassService()
			password, err := testService.GeneratePassword(testCase.expectedLength)

			assert.ErrorIs(t, err, testCase.expectedError)
			assert.Equal(t, testCase.expectedLength, len(password))
		})
	}
}
