package psd

import (
	"io"

	"github.com/yu-ichiko/go-psd/section/colormodedata"
	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/section/layer"
	"github.com/yu-ichiko/go-psd/section/resources"
)

type PSD struct {
	Header         Header
	ColorModeData  ColorModeData
	ImageResources []ImageResourcesBlock
}

func Parse(r io.Reader) error {

	header, _, err := header.Parse(r)
	if err != nil {
		return err
	}

	//fmt.Printf("%+v\n", header)

	_, _, err = colormodedata.Parse(r, header)
	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", colorModeData)

	_, _, err = resources.Parse(r)
	if err != nil {
		return err
	}

	// fmt.Println(imageResourceBlocks)

	_, _, err = layer.Parse(r, header)
	if err != nil {
		return err
	}

	return nil
}
