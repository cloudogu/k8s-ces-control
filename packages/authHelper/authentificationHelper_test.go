package authHelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-ces-control/packages/account"
)

func Test_GetServiceAccountCredentials(t *testing.T) {
	t.Run("should not check authHelper as method belongs to public api. No authentication required", func(t *testing.T) {
		credentials, err := GetServiceAccountCredentials(context.TODO(), "testservice")
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create ServiceAccountManager: failed to get key provider from global config")
		assert.Equal(t, credentials, account.ServiceAccountData{})
	})
}
