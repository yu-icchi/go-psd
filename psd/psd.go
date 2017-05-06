package psd

import (
	"io"
	"fmt"

	"github.com/yu-ichiko/go-psd/psd/section/header"
	"github.com/yu-ichiko/go-psd/psd/section/colormodedata"
	"github.com/yu-ichiko/go-psd/psd/section/resources"
	"github.com/yu-ichiko/go-psd/psd/section/layer"
)

func Parse(r io.Reader) error {

	header, _, err := header.Parse(r)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", header)

	colorModeData, _, err := colormodedata.Parse(r)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", colorModeData)

	imageResourceBlocks, _, err := resources.Parse(r)
	if err != nil {
		return err
	}

	fmt.Println(imageResourceBlocks)

	layer.Parse(r, header)

	return nil
}
