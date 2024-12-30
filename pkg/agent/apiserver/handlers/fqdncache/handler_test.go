package fqdncache

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"antrea.io/antrea/pkg/agent/apis"
	queriertest "antrea.io/antrea/pkg/agent/querier/testing"
	"antrea.io/antrea/pkg/agent/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFqdnCacheQuery(t *testing.T) {
	tests := []struct {
		name             string
		expectedStatus   int
		expectedResponse []types.DnsCacheEntry
	}{
		{
			name:           "FQDN cache exists",
			expectedStatus: http.StatusOK,
			expectedResponse: []types.DnsCacheEntry{
				{
					FqdnName:       "google.com",
					IpAddress:      net.ParseIP("10.0.0.1"),
					ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
				},
			},
		},
		{
			name:           "FQDN cache does not exist",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			q := queriertest.NewMockAgentQuerier(ctrl)
			q.EXPECT().GetFqdnCache().Return(tt.expectedResponse[0])
			handler := HandleFunc(q)

			req, err := http.NewRequest(http.MethodGet, "", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedStatus == http.StatusOK {
				var received []apis.FqdnCacheResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &received)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, received)
			}
		})
	}
}
