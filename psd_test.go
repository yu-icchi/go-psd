package psd

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	file, err := os.Open("./testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	err = Parse(file)
	assert.NoError(t, err)
}
