package psd

import (
	"io"

	"github.com/yu-ichiko/go-psd/section/colormodedata"
	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/section/layer"
	"github.com/yu-ichiko/go-psd/section/resources"
)

type PSD struct {
	Header         *Header
	ColorModeData  *ColorModeData
	ImageResources []*ImageResourceBlock
	Layers         []*Layer
}

const (
	sectionLen     = 4
	compressionLen = 2
	sigLen         = 4
)

func newPixel(h *Header, hasAlpha bool) {

}
