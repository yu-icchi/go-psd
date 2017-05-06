package layer

import (
	"io"

	"fmt"
	"github.com/yu-ichiko/go-psd/psd/header"
	"github.com/yu-ichiko/go-psd/psd/util"
)

func parseChannelImageData(r io.Reader, header *header.Header, layer *Layer) (read int, err error) {
	var l int

	//imageChs := make([][]byte, header.ColorMode.Channels(), 8)
	//fmt.Println(imageChs)

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

		switch compMethod {
		case 0:
			parseRaw(r, layer)
		case 1:
			parseRLE(r, header, layer)
		case 2, 3:
			parseZip()
		}
	}

	return
}

func parseRaw(r io.Reader, layer *Layer) error {
	buf := make([]byte, layer.Width() * layer.Height())
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}

	fmt.Println("parseRaw len:", len(buf))

	return nil
}

func parseRLE(r io.Reader, header *header.Header, layer *Layer) (read int, err error) {
	var l int

	lines := layer.Height()
	//fmt.Println("== lines:", lines)
	//fmt.Println("== ", lines*(util.GetSize(header.IsPSB())>>1))
	buf := make([]byte, lines*(util.GetSize(header.IsPSB())>>1))
	if l, err = io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	read += l

	//fmt.Println(buf)

	total := 0
	lens := make([]int, lines)
	offsets := make([]int, lines)
	ofs := 0
	if header.IsPSB() {
		for i := range lens {
			l = int(util.ReadUint32(buf, ofs))
			lens[i] = l
			offsets[i] = total
			total += l
			ofs += 4
		}
	} else {
		for i := range lens {
			l = int(util.ReadUint16(buf, ofs))
			lens[i] = l
			offsets[i] = total
			total += l
			ofs += 2
		}
	}

	//fmt.Println(total)

	buf = make([]byte, total)
	if l, err = io.ReadFull(r, buf); err != nil {
		fmt.Println(err)
		return
	}
	read += l

	//fmt.Println(buf)

	return read, nil
}

func parseZip() {

}
