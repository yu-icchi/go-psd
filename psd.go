package psd

import (
	"io"
)

func Decode(r io.Reader) (*PSD, error) {

	psd := newPSD()
	var read int

	// File Header Section
	// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_19840
	header, l, err := readHeader(r)
	if err != nil {
		return nil, err
	}
	psd.SetFileHeader(*header)
	read += l

	// Color Mode Data Section
	// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_71638
	colorModeData, l, err := readColorModeData(r)
	if err != nil {
		return nil, err
	}
	psd.SetColorModeData(colorModeData)
	read += l

	// Image Resources Section
	// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_69883
	imageRes, _, err := readImageResources(r)
	if err != nil {
		return nil, err
	}
	psd.ImageResourceBlocks = imageRes

	err = readLayerAndMarkInfo(r, psd)
	if err != nil {
		return nil, err
	}

	return psd, nil
}

func Encode(w io.Writer, psd *PSD) error {
	return nil
}
