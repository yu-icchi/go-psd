package psd

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	file, err := os.Open("./testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	err = Parse(file)
	assert.NoError(t, err)
}

func BenchmarkParse(b *testing.B) {
	file, err := os.Open("./testdata/test.psd")
	if err != nil {
		b.FailNow()
	}
	defer file.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = file.Seek(0, 0)
		if err != nil {
			b.Fatal(err)
		}
		err = Parse(file)
		if err != nil {
			b.Fatal(err)
		}
	}
}
