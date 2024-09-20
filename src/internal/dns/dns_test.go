// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		serviceURL        string
		expectedErr       string
		expectedNamespace string
		expectedName      string
		expectedPort      int
	}{
		{
			name:              "correct service url",
			serviceURL:        "http://foo.bar.svc.cluster.local:5000",
			expectedNamespace: "bar",
			expectedName:      "foo",
			expectedPort:      5000,
		},
		{
			name:        "invalid service url without port",
			serviceURL:  "http://google.com",
			expectedErr: "service url does not have a port",
		},
		{
			name:        "invalid service url with port",
			serviceURL:  "http://google.com:3000",
			expectedErr: "invalid service url http://google.com:3000",
		},
		{
			name:        "empty service url",
			serviceURL:  "",
			expectedErr: "service url cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isServiceURL := IsServiceURL(tt.serviceURL)
			namespace, name, port, err := ParseServiceURL(tt.serviceURL)
			if tt.expectedErr != "" {
				require.False(t, isServiceURL)
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.True(t, isServiceURL)
			require.Equal(t, tt.expectedNamespace, namespace)
			require.Equal(t, tt.expectedName, name)
			require.Equal(t, tt.expectedPort, port)
		})
	}
}