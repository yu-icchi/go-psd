package header

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestParse(t *testing.T) {
	file, err := os.Open("../../testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	header, read, err := Parse(file)
	assert.NoError(t, err)
	fmt.Println(header)
	fmt.Println(read)
}
