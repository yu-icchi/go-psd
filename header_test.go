package psd

import (
	"testing"
	"os"
	"fmt"
	"bufio"
	"io"
)

func Test_readHeader(t *testing.T) {
	file, err := os.Open("testdata/test_illust_photoshop_1027.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()
	header, err := readHeader(file)
	fmt.Println(header)
	fmt.Println(err)

	fmt.Println(header.ColorMode())

	wf, err := os.Create("testdata/sample.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	writer := bufio.NewWriter(wf)
	writeHeader(writer, header)

	file, err = os.Open("testdata/sample.psd")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	buf := make([]byte, 26)
	if _, err := io.ReadFull(file, buf); err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(buf)
}
