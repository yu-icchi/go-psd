package enginedata

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestDecode(t *testing.T) {
	file, err := ioutil.ReadFile("./testdata/enginedata")
	require.NoError(t, err)

	data, err := Decode(file)
	assert.NoError(t, err)
	fmt.Println(data)
}
