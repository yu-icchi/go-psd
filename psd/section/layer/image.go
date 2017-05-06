package layer

import (
	"fmt"
	"io"

	"github.com/yu-ichiko/go-psd/psd/section/header"
	"github.com/yu-ichiko/go-psd/psd/util"
)

func parseChannelImageData(r io.Reader, header *header.Header, layer *Layer) (read int, err error) {
	var l int

	img := map[int][]byte{}
	for _, channel := range layer.Channels {

		buf := make([]byte, 2)
		if l, err = io.ReadFull(r, buf); err != nil {
			return
		}
		read += l
		compMethod := int(util.ReadUint16(buf, 0))

		if channel.Length == 2 {
			continue
		}

		fmt.Println("++ compMethod:", compMethod)

		var imgCh []byte
		switch compMethod {
		case 0:
			imgCh, l, err = parseRaw(r, layer)
		case 1:
			imgCh, l, err = parseRLE(r, header, layer)
		case 2, 3:
			parseZip()
		default:
			return read, fmt.Errorf("unknown compression method: %d", compMethod)
		}

		img[channel.ID] = imgCh
	}

	fmt.Println(len(img))

	return
}

func parseRaw(r io.Reader, layer *Layer) (img []byte, read int, err error) {
	var l int
	img = make([]byte, layer.Width()*layer.Height())
	if l, err = io.ReadFull(r, img); err != nil {
		return nil, read, err
	}
	read += l

	return img, read, nil
}

func parseRLE(r io.Reader, header *header.Header, layer *Layer) (img []byte, read int, err error) {
	var l int
	buf := make([]byte, layer.Height()*(util.GetSize(header.IsPSB())/2))
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	read += l

	var total, n int
	for i := 0; i < layer.Height(); i++ {
		if header.IsPSB() {
			l = int(util.ReadUint32(buf, n))
			n += 4
		} else {
			l = int(util.ReadUint16(buf, n))
			n += 2
		}
		total += l
	}

	img = make([]byte, total)
	if l, err = io.ReadFull(r, img); err != nil {
		return nil, read, err
	}
	read += l

	return img, read, nil
}

func parseZip() {

}
