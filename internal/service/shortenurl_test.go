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

const (
	testURL    = "http://localhost:8000"
	testExpire = time.Hour
)

var errCodeNotFound = errors.New("code not found")

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
				mockRepo.On("GetUrlByCode", context.Background(), "1234567").Return("", errCodeNotFound)
				mockRepo.On("StoreUrl", context.Background(), "1234567", testURL, testExpire).Return(nil)
				return mockRepo
			},
			expectedErr:  nil,
			expectedResp: "1234567",
		},

		{
			name: "retry when code already exists",
			setUpMockGenCode: func(t *testing.T) *mocksGenCode.GenPass {
				mockGenCode := mocksGenCode.NewGenPass(t)
				mockGenCode.On("GeneratePassword", 7).Return("1111111", nil).Once()
				mockGenCode.On("GeneratePassword", 7).Return("2222222", nil).Once()
				return mockGenCode
			},
			setUpMockRepo: func(t *testing.T) *mocksUrlStorage.UrlStorage {
				mockRepo := mocksUrlStorage.NewUrlStorage(t)
				mockRepo.On("GetUrlByCode", context.Background(), "1111111").Return(testURL, nil).Once()
				mockRepo.On("GetUrlByCode", context.Background(), "2222222").Return("", errCodeNotFound).Once()
				mockRepo.On("StoreUrl", context.Background(), "2222222", testURL, testExpire).Return(nil).Once()
				return mockRepo
			},
			expectedErr:  nil,
			expectedResp: "2222222",
		},

		{
			name: "error generate password",
			setUpMockGenCode: func(t *testing.T) *mocksGenCode.GenPass {
				mockGenCode := mocksGenCode.NewGenPass(t)
				// Không quan tâm length cụ thể vì fail ngay ở bước generate
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
				mockRepo.On("GetUrlByCode", context.Background(), "1234567").Return("", errCodeNotFound)
				mockRepo.On("StoreUrl", context.Background(), "1234567", testURL, testExpire).Return(errors.New("Internal server error"))
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

			shortenUrlService := NewShortenUrl(mockRepo, mockGenCode)

			shortenUrl, err := shortenUrlService.CreateShortenUrl(context.Background(), testURL, testExpire)

			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedResp, shortenUrl)
		})
	}
}

func TestGetShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setUpMockRepo func(t *testing.T) *mocksUrlStorage.UrlStorage
		ExpectErr     error
		ExpectUrl     string
	}{
		{
			name: "success",
			setUpMockRepo: func(t *testing.T) *mocksUrlStorage.UrlStorage {
				mockRepo := mocksUrlStorage.NewUrlStorage(t)
				mockRepo.On("GetUrlByCode", context.Background(), "123456").Return(testURL, nil)
				return mockRepo
			},
			ExpectErr: nil,
			ExpectUrl: testURL,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := testCase.setUpMockRepo(t)

			shortenUrlService := NewShortenUrl(mockRepo, nil)

			url, err := shortenUrlService.GetUrlByCode(context.Background(), "123456")

			assert.Equal(t, testCase.ExpectErr, err)
			assert.Equal(t, testCase.ExpectUrl, url)
		})
	}
}
