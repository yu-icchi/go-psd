package colormodedata

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

	file.Seek(26, 0)

	data, read, err := Parse(file)
	assert.NoError(t, err)
	fmt.Println(data)
	fmt.Println(read)
}