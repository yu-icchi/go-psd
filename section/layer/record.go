package layer

import (
	"errors"
	"io"

	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

func parseRecord(r io.Reader, header *header.Header) (layer Layer, read int, err error) {
	var l int

	// Rectangle containing the contents of the layer
	buf := make([]byte, 4*4+2)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	read += l

	layer = Layer{}

	layer.Top = int(util.ReadUint32(buf, 0))
	layer.Left = int(util.ReadUint32(buf, 4))
	layer.Bottom = int(util.ReadUint32(buf, 8))
	layer.Right = int(util.ReadUint32(buf, 12))

	// Number of channels in the layer
	numChannels := int(util.ReadUint16(buf, 16))

	size := util.GetSize(header.IsPSB())
	channels := make([]Channel, numChannels)
	for i := range channels {
		channel := Channel{}

		buf := make([]byte, 2+size)
		if l, err = io.ReadFull(r, buf); err != nil {
			return
		}
		read += l

		channel.ID = int(int16(util.ReadUint16(buf, 0)))
		if size == 4 {
			channel.Length = int(util.ReadUint32(buf, 2))
		}
		if size == 8 {
			channel.Length = int(util.ReadUint64(buf, 2))
		}

		channels[i] = channel
	}
	layer.Channels = channels

	buf = make([]byte, 12)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	read += l

	if util.ReadString(buf, 0, 4) != "8BIM" {
		err = errors.New("psd: invalid layer signature")
		return
	}

	// Blend Mode
	layer.BlendMode = util.ReadString(buf, 4, 8)
	// Opacity
	layer.Opacity = int(buf[8])
	// Clipping
	layer.Clipping = Clipping(buf[9])
	// Flags
	layer.Flags = buf[10]
	// Filler
	layer.Filter = int(buf[11])

	//fmt.Printf("%b\n", layer.Flags)
	//fmt.Println("+++", layer.Visible())
	//fmt.Println("+++", layer.Obsolete())

	// extra field length
	buf = make([]byte, 4)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	extraLength := int(util.ReadUint32(buf, 0))
	if extraLength <= 0 {
		return
	}

	buf = make([]byte, extraLength)
	if l, err = io.ReadFull(r, buf); err != nil {
		return
	}
	read += l

	// Layer mask / adjustment layer data
	n, err := parseMask(buf)

	// Layer blending ranges data
	m, err := parseBlendingRangesData(buf[n:])

	// Layer name (MBCS)
	str, l := util.PascalString(buf, n+m)
	layer.LegacyName = str
	//padding := (4 - ((1 + l) % 4)) % 4
	//fmt.Println(padding)

	return
}
