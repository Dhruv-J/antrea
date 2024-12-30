package fqdncache

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"
	"time"

	"antrea.io/antrea/pkg/agent/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestTrasnform(t *testing.T) {
	var fqdn1 = types.DnsCacheEntry{
		FqdnName:       "google.com",
		IpAddress:      net.ParseIP("10.0.0.1"),
		ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
	}
	var fqdn2 = types.DnsCacheEntry{
		FqdnName:       "google.com",
		IpAddress:      net.ParseIP("10.0.0.2"),
		ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
	}
	var fqdn3 = types.DnsCacheEntry{
		FqdnName:       "google.com",
		IpAddress:      net.ParseIP("10.0.0.3"),
		ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
	}
	var fqdn4 = types.DnsCacheEntry{
		FqdnName:       "example.com",
		IpAddress:      net.ParseIP("10.0.0.4"),
		ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
	}
	var fqdn5 = types.DnsCacheEntry{
		FqdnName:       "antrea.io",
		IpAddress:      net.ParseIP("10.0.0.5"),
		ExpirationTime: time.Date(2025, 12, 25, 15, 0, 0, 0, time.Now().Location()),
	}
	var fqdnList = []types.DnsCacheEntry{fqdn1, fqdn2, fqdn3, fqdn4, fqdn5}

	tests := []struct {
		name             string
		opts             map[string]string
		fqdnList         []types.DnsCacheEntry
		expectedResponse interface{}
		expectedError    string
	}{
		{
			name:             "all",
			fqdnList:         fqdnList,
			expectedResponse: []Response{{&fqdn1}, {&fqdn2}, {&fqdn3}, {&fqdn4}, {&fqdn5}},
		},
		{
			name: "only google.com domain name",
			opts: map[string]string{
				"domain": "google.com",
			},
			fqdnList:         fqdnList,
			expectedResponse: []Response{{&fqdn1}, {&fqdn2}, {&fqdn3}},
		},
		{
			name: "only antrea.io domain name",
			opts: map[string]string{
				"domain": "antrea.io",
			},
			fqdnList:         fqdnList,
			expectedResponse: []Response{{&fqdn5}},
		},
		{
			name: "domain name that doesn't exist",
			opts: map[string]string{
				"domain": "bing.com",
			},
			fqdnList:         fqdnList,
			expectedResponse: []Response{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqByte, _ := json.Marshal(tt.fqdnList)
			reqReader := bytes.NewReader(reqByte)
			result, err := Transform(reqReader, false, tt.opts)
			if tt.expectedError == "" {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, result)
			} else {
				assert.ErrorContains(t, err, tt.expectedError)
			}
		})
	}
}
