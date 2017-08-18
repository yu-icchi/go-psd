package additional

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLayerID(t *testing.T) {
	buf := []byte{0, 0, 2, 214}
	id, err := NewLayerID(buf)
	require.NoError(t, err)
	assert.Equal(t, 726, id)
}
