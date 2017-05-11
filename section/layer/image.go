package layer

import (
	"fmt"
	"image"
	"io"

	"github.com/yu-ichiko/go-psd/pixel"
	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

func parseChannelImageData(r io.Reader, header *header.Header, layer Layer) (image.Image, int, error) {
	var l, read int
	var err error

	img := map[int][]byte{}
	for _, channel := range layer.Channels {

		buf := make([]byte, 2)
		if l, err = io.ReadFull(r, buf); err != nil {
			return nil, read, err
		}
		read += l
		compMethod := int(util.ReadUint16(buf, 0))

		if channel.Length == 2 {
			continue
		}

		var imgCh []byte
		switch compMethod {
		case 0:
			imgCh, l, err = parseRaw(r, layer)
		case 1:
			imgCh, l, err = parseRLE(r, header, layer)
		case 2, 3:
			parseZip()
		default:
			return nil, read, fmt.Errorf("unknown compression method: %d", compMethod)
		}

		img[channel.ID] = imgCh
	}

	if len(img) <= 0 {
		return nil, read, nil
	}

	p := pixel.NewPixel(header, header.ColorMode.Channels() < len(layer.Channels))
	p.SetSource(layer.Top, layer.Left, layer.Bottom, layer.Right, img[0], img[1], img[2], img[-1])

	return p, read, nil
}

func parseRaw(r io.Reader, layer Layer) (img []byte, read int, err error) {
	var l int
	img = make([]byte, layer.Width()*layer.Height())
	if l, err = io.ReadFull(r, img); err != nil {
		return nil, read, err
	}
	read += l

	return img, read, nil
}

func parseRLE(r io.Reader, header *header.Header, layer Layer) (img []byte, read int, err error) {
	var l int
	buf := make([]byte, layer.Height()*(util.GetSize(header.IsPSB())/2))
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}
	read += l

	lens := make([]int, layer.Height())
	var total, n int
	for i := 0; i < layer.Height(); i++ {
		if header.IsPSB() {
			l = int(util.ReadUint32(buf, n))
			n += 4
		} else {
			l = int(util.ReadUint16(buf, n))
			n += 2
		}
		lens[i] = l
		total += l
	}

	buf = make([]byte, total)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l

	dest := make([]byte, (layer.Width()*header.Depth+7)>>3*layer.Height())
	decodePackBitsPerLine(dest, buf, lens)

	return dest, read, nil
}

func decodePackBitsPerLine(dest []byte, buf []byte, lens []int) {
	var l int
	for _, ln := range lens {
		for i := 0; i < ln; {
			if buf[i] <= 0x7f {
				l = int(buf[i]) + 1
				copy(dest[:l], buf[i+1:])
				dest = dest[l:]
				i += l + 1
				continue
			}
			l = int(-buf[i]) + 1
			for j, c := 0, buf[i+1]; j < l; j++ {
				dest[j] = c
			}
			dest = dest[l:]
			i += 2
		}
		buf = buf[ln:]
	}
}

func parseZip() {

}
