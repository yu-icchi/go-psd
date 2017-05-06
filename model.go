package psd

import "image"

type PSD struct {
	FileHeader          FileHeader
	ColorModeData       ColorModeData
	ImageResourceBlocks []ImageResourceBlock
}

type FileHeader struct {
	Signature string
	Version   int
	Channels  int
	Height    int
	Width     int
	Depth     int
	ColorMode int
}

type ColorModeData []byte

type ImageResourceBlock struct {
	Signature string
	ID        int
	Name      string
	Data      []byte
}

type Layer struct {
	ID          int
	Index       int
	Name        string
	X           int
	Y           int
	Width       int
	Height      int
	Opacity     uint8
	ColorMode   string
	Clipping    byte
	Flags       byte
	ChannelInfo []ChannelInfo
	BlendMode   BlendMode
	Image       image.Image
}

type ChannelInfo struct {
	ID   int
	Data uint64
}

type BlendMode struct {
	Signature string
	Key       string
}

type AddLayerInfo struct {
	Signature string
	Key       string
	Data      []byte
}

func newPSD() *PSD {
	return &PSD{}
}

func (psd *PSD) IsPSB() bool {
	return psd.FileHeader.Version == 2
}

func (psd *PSD) SetFileHeader(header FileHeader) {
	psd.FileHeader = header
}

func (psd *PSD) SetColorModeData(data []byte) {
	psd.ColorModeData = data
}
