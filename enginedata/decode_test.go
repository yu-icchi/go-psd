package enginedata

import (
	"encoding/json"
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
	jsonByte, err := json.Marshal(data)
	assert.NoError(t, err)
	fmt.Println(string(jsonByte))
}

func BenchmarkDecode(b *testing.B) {
	file, _ := ioutil.ReadFile("./testdata/enginedata")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(file)
	}
}
