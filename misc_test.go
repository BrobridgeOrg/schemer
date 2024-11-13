package schemer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMisc_ParseEntryPath(t *testing.T) {

	key, index := parsePathEntry("a.b[888]")
	assert.Equal(t, "a.b", key)
	assert.Equal(t, 888, index)
}
