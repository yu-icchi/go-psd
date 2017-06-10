package layer

import (
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

func parseInfo(r io.Reader, header *header.Header) ([]Layer, int, error) {
	var l, read int
	var err error

	// Length of the layers info section
	size := util.GetSize(header.IsPSB())
	buf := make([]byte, size)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	read += l
	size = int(util.ReadUint(buf))
	if size <= 0 {
		return nil, read, nil
	}

	// Layer count
	if l, err = io.ReadFull(r, buf[:2]); err != nil {
		return nil, read, err
	}
	read += l
	count := util.Abs(int(util.ReadInt16(buf, 0)))
	// fmt.Println("=== info count:", count)

	layers := make([]Layer, count)
	var i int
	for i = 0; i < count; i++ {
		layer, l, err := parseRecord(r, header)
		if err != nil {
			return nil, read, err
		}
		read += l
		layer.Index = i
		//fmt.Printf("==== record layer: %+v\n", layer)
		layers[i] = layer
	}

	// Channel image data
	for i = 0; i < count; i++ {
		img, l, err := parseChannelImageData(r, header, layers[i])
		if err != nil {
			return nil, read, err
		}
		read += l
		layers[i].Image = img
	}

	return layers, read, nil
}
