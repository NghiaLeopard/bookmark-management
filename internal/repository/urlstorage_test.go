package repository

import (
	"context"
	"testing"
	"time"

	rdb "github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestStoreUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		SetUpMockRdb func(t *testing.T) *redis.Client
		ExpectErr    error
		VerifyFunc   func(t *testing.T, rdb *redis.Client)
	}{
		{
			name: "Normal case",
			SetUpMockRdb: func(t *testing.T) *redis.Client {
				return rdb.NewMockRClient(t)
			},
			ExpectErr: nil,
			VerifyFunc: func(t *testing.T, rdb *redis.Client) {
				res, err := rdb.Get(context.Background(), "test").Result()
				assert.NoError(t, err)
				assert.Equal(t, "http://localhost:8000", res)
			},
		},

		{
			name: "Connection error",
			SetUpMockRdb: func(t *testing.T) *redis.Client {
				mockRdb := rdb.NewMockRClient(t)

				mockRdb.Close()
				return mockRdb
			},
			ExpectErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRdb := tc.SetUpMockRdb(t)

			rdb := NewUrlStorage(mockRdb)
			err := rdb.StoreUrl(context.Background(), "test", "http://localhost:8000", time.Hour)
			assert.Equal(t, tc.ExpectErr, err)

			if tc.VerifyFunc != nil {
				tc.VerifyFunc(t, mockRdb)
			}
		})
	}
}

func TestGetUrlByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		SetUpMockRdb  func(t *testing.T) *redis.Client
		SetupStoreUrl func(t *testing.T, rdb UrlStorage)
		ExpectErr     error
		ExpectUrl     string
	}{
		{
			name: "Normal case",
			SetUpMockRdb: func(t *testing.T) *redis.Client {
				return rdb.NewMockRClient(t)
			},
			SetupStoreUrl: func(t *testing.T, rdb UrlStorage) {
				err := rdb.StoreUrl(context.Background(), "123456", "http://localhost:8000", time.Hour)
				assert.NoError(t, err)
			},
			ExpectErr: nil,
			ExpectUrl: "http://localhost:8000",
		},
		{
			name: "Code not found",
			SetUpMockRdb: func(t *testing.T) *redis.Client {
				mockRdb := rdb.NewMockRClient(t)
				return mockRdb
			},
			SetupStoreUrl: nil,
			ExpectErr:     ErrCodeNotFound,
			ExpectUrl:     "",
		},
		{
			name: "Connection error",
			SetUpMockRdb: func(t *testing.T) *redis.Client {
				mockRdb := rdb.NewMockRClient(t)
				mockRdb.Close()
				return mockRdb
			},
			SetupStoreUrl: nil,
			ExpectErr:     redis.ErrClosed,
			ExpectUrl:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRdb := tc.SetUpMockRdb(t)

			rdb := NewUrlStorage(mockRdb)

			if tc.SetupStoreUrl != nil {
				tc.SetupStoreUrl(t, rdb)
			}

			url, err := rdb.GetUrlByCode(context.Background(), "123456")
			assert.Equal(t, tc.ExpectErr, err)
			assert.Equal(t, tc.ExpectUrl, url)
		})
	}

}
