package layer

import (
	"errors"
	"fmt"
	"io"

	"github.com/yu-ichiko/go-psd/psd/section/header"
	"github.com/yu-ichiko/go-psd/psd/util"
)

func parseRecord(r io.Reader, header *header.Header) (*Layer, int, error) {
	var l int
	var read int
	var err error

	// Rectangle containing the contents of the layer
	buf := make([]byte, 4*4)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l

	layer := Layer{}

	layer.Top = int(util.ReadUint32(buf, 0))
	layer.Left = int(util.ReadUint32(buf, 4))
	layer.Bottom = int(util.ReadUint32(buf, 8))
	layer.Right = int(util.ReadUint32(buf, 12))

	// Number of channels in the layer
	buf = make([]byte, 2)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l
	numChannels := int(util.ReadUint16(buf, 0))

	channels := []Channel{}
	for i := 0; i < numChannels; i++ {
		channel := Channel{}

		buf := make([]byte, 2)
		if l, err = io.ReadFull(r, buf); err != nil {
			return nil, read, err
		}
		read += l
		channel.ID = int(int16(util.ReadUint16(buf, 0)))

		size := util.GetSize(header.IsPSB())
		buf = make([]byte, size)
		if l, err = io.ReadFull(r, buf); err != nil {
			return nil, read, err
		}
		read += l
		channel.Length = int(util.ReadUint(buf))

		channels = append(channels, channel)
	}
	layer.Channels = channels

	buf = make([]byte, 12)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l

	if util.ReadString(buf, 0, 4) != "8BIM" {
		return nil, read, errors.New("invalid psd:layer signature")
	}

	// Blend Mode
	layer.BlendMode = util.ReadString(buf, 4, 8)
	// Opacity
	layer.Opacity = int(buf[8])
	// Clipping
	layer.Clipping = int(buf[9])
	// Flags
	layer.Flags = int(buf[10])
	// Filler
	layer.Filter = int(buf[11])

	// extra field length
	buf = make([]byte, 4)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	extraLength := int(util.ReadUint32(buf, 0))
	if extraLength <= 0 {
		return nil, read, err
	}

	buf = make([]byte, extraLength)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l

	// Layer mask / adjustment layer data
	n, err := parseMask(buf)

	// Layer blending ranges data
	m, err := parseBlendingRangesData(buf[n:])

	// Layer name (MBCS)
	str, l := util.PascalString(buf, n+m)
	layer.LegacyName = str
	padding := (4 - ((1 + l) % 4)) % 4
	fmt.Println(padding)

	fmt.Println(buf[n+m+l:])

	return &layer, read, nil
}
