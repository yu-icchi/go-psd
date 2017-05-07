package psd

import (
	"io"
	"fmt"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/section/colormodedata"
	"github.com/yu-ichiko/go-psd/section/resources"
	"github.com/yu-ichiko/go-psd/section/layer"
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

	layers, _, err := layer.Parse(r, header)
	if err != nil {
		return err
	}
	for _, l := range layers {
		if l.Image != nil {
			fmt.Println(l.Image.Bounds())
		}
	}

	return nil
}
