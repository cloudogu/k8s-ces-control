package authHelper

import (
	"context"
	"github.com/cloudogu/k8s-ces-control/packages/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetServiceAccountCredentials(t *testing.T) {
	t.Run("should not check authHelper as method belongs to public api. No authentication required", func(t *testing.T) {
		credentials, err := GetServiceAccountCredentials(context.TODO(), "testservice")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find configuration")
		assert.Equal(t, credentials, account.ServiceAccountData{})
	})
}
