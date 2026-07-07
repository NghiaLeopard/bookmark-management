package service

import (
	"context"
	"errors"
	"testing"
	"time"

	mocksUrlStorage "github.com/NghiaLeopard/bookmark-management/internal/repository/mocks"
	mocksGenCode "github.com/NghiaLeopard/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setUpMockRepo    func(t *testing.T) *mocksUrlStorage.UrlStorage
		setUpMockGenCode func(t *testing.T) *mocksGenCode.GenPass
		expectedErr      error
		expectedResp     string
	}{
		{
			name: "success",
			setUpMockGenCode: func(t *testing.T) *mocksGenCode.GenPass {
				mockGenCode := mocksGenCode.NewGenPass(t)
				mockGenCode.On("GeneratePassword", 7).Return("1234567", nil)
				return mockGenCode
			},
			setUpMockRepo: func(t *testing.T) *mocksUrlStorage.UrlStorage {
				mockRepo := mocksUrlStorage.NewUrlStorage(t)
				mockRepo.On("StoreUrl", context.Background(), "1234567", "http://localhost:8000", time.Hour).Return(nil)
				return mockRepo
			},
			expectedErr:  nil,
			expectedResp: "1234567",
		},

		{
			name: "error generate password",
			setUpMockGenCode: func(t *testing.T) *mocksGenCode.GenPass {
				mockGenCode := mocksGenCode.NewGenPass(t)
				mockGenCode.On("GeneratePassword", mock.Anything).Return("", errors.New("Internal server error"))
				return mockGenCode
			},
			setUpMockRepo: func(t *testing.T) *mocksUrlStorage.UrlStorage {
				mockRepo := mocksUrlStorage.NewUrlStorage(t)
				return mockRepo
			},
			expectedErr:  errors.New("Internal server error"),
			expectedResp: "",
		},

		{
			name: "error store url",
			setUpMockRepo: func(t *testing.T) *mocksUrlStorage.UrlStorage {
				mockRepo := mocksUrlStorage.NewUrlStorage(t)
				mockRepo.On("StoreUrl", context.Background(), "1234567", "http://localhost:8000", time.Hour).Return(errors.New("Internal server error"))
				return mockRepo
			},
			setUpMockGenCode: func(t *testing.T) *mocksGenCode.GenPass {
				mockGenCode := mocksGenCode.NewGenPass(t)
				mockGenCode.On("GeneratePassword", 7).Return("1234567", nil)
				return mockGenCode
			},
			expectedErr:  errors.New("Internal server error"),
			expectedResp: "",
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := testCase.setUpMockRepo(t)

			mockGenCode := testCase.setUpMockGenCode(t)

			shortenUrlService := NewShortenUrlService(mockRepo, mockGenCode)

			shortenUrl, err := shortenUrlService.CreateShortenUrl(context.Background(), "http://localhost:8000", time.Hour)

			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedResp, shortenUrl)
		})
	}
}
