package repository

import (
	"context"
	"testing"
	"time"

	rdb "github.com/NghiaLeopard/bookmark-management/internal/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestUrlStorage(t *testing.T) {
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
