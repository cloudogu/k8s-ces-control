package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_randomString(t *testing.T) {
	t.Run("should generate string with length 8", func(t *testing.T) {
		generatedString := randomString(8)

		assert.Len(t, generatedString, 8)
	})
	t.Run("should generate nothing if length parameter is zero", func(t *testing.T) {
		generatedString := randomString(0)

		assert.Len(t, generatedString, 0)
	})
	t.Run("should generate nothing if length parameter is negative", func(t *testing.T) {
		generatedString := randomString(-10)

		assert.Len(t, generatedString, 0)
	})
}
