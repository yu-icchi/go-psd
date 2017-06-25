package psd

import "image"

const (
	sectionLen     = 4
	compressionLen = 2
	sigLen         = 4

	imgRAW                  = 0
	imgRLE                  = 1
	imgZIPWithOutPrediction = 2
	imgZIPWithPrediction    = 3
)

type PSD struct {
	Header          *Header
	ColorModeData   *ColorModeData
	ImageResources  []*ImageResourceBlock
	Layers          []*Layer
	GlobalLayerMask *GlobalLayerMask
	AdditionalInfos []*AdditionalInfo
	Image           image.Image
}
