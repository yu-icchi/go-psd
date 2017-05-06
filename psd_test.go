package psd

import (
	"testing"
	"os"
	"fmt"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestDecode(t *testing.T) {
	file, err := os.Open("testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	psd, err := Decode(file)
	assert.NoError(t, err)
	fmt.Println(psd)

	jsonByte, err := json.Marshal(psd)
	assert.NoError(t, err)
	fmt.Println(string(jsonByte))
}
