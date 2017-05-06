package psd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDecode_readHeader(t *testing.T) {
	file, err := os.Open("testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	header, read, err := readHeader(file)
	assert.NoError(t, err)
	fmt.Println(header)
	fmt.Println(read)
}
