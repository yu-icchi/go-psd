package additional

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLayerName(t *testing.T) {
	buf := []byte{0, 0, 0, 6, 48, 236, 48, 164, 48, 228, 48, 252, 0, 32, 0, 48}
	name, err := NewLayerName(buf)
	require.NoError(t, err)
	assert.Equal(t, "レイヤー 0", name)
}
