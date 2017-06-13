package psd

import "image"

type PSD struct {
	Header         *Header
	ColorModeData  *ColorModeData
	ImageResources []*ImageResourceBlock
	Layers         []*Layer
	Image          image.Image
}

const (
	sectionLen     = 4
	compressionLen = 2
	sigLen         = 4

	imgRAW                  = 0
	imgRLE                  = 1
	imgZIPWithoutPrediction = 2
	imgZIPWithPrediction    = 3
)
