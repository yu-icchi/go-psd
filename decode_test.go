package psd

import (
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"strconv"
	"testing"
)

func TestDecode(t *testing.T) {
	filename := "mask"
	file, err := os.Open("./testdata/" + filename + ".psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	psd, err := Decode(file)
	require.NoError(t, err)

	if psd.Image != nil {
		savePNG("./png/"+filename+".png", psd.Image)
	}
	for _, layer := range psd.Layers {
		if layer.Image == nil {
			continue
		}
		pp.Println(layer.AdditionalInfos)

		filename := "./png/" + strconv.Itoa(layer.ID) + ".png"
		if err := savePNG(filename, layer.Image); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	//output, err := os.Create("./testdata/" + filename + "_enc.psd")
	//if err != nil {
	//	t.Error(err)
	//	t.FailNow()
	//}
	//defer output.Close()

	//err = Encode(output, psd)
	//require.NoError(t, err)
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

// 943,868,237
// 943868237
