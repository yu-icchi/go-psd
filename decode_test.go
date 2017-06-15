package psd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"os"
	"strconv"
	"testing"
)

func TestDecode(t *testing.T) {
	file, err := os.Open("./testdata/test.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	psd, err := Decode(file)
	assert.NoError(t, err)

	// savePNG("./png/test.png", psd.Image)

	fmt.Println(psd.Header)
	for _, layer := range psd.Layers {
		if layer.Image == nil {
			continue
		}

		filename := "./png/" + strconv.Itoa(layer.Index) + ".png"
		if err := savePNG(filename, layer.Image); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

func savePNG(name string, img image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func BenchmarkDecode(b *testing.B) {
	file, err := os.Open("./testdata/test.psd")
	if err != nil {
		b.Error(err)
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
		Decode(file)
	}
}
