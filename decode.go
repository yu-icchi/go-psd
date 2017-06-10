package psd

import (
	"bytes"
	"errors"
	"image"
	"io"

	"fmt"
	"github.com/yu-ichiko/go-psd/pixel"
	"github.com/yu-ichiko/go-psd/util"
)

type decoder struct {
	r    io.Reader
	buf  []byte
	read int

	header *Header
	sections
}

type sections struct {
	colorModeData           int
	imageResources          int
	layerAndMaskInformation int
}

func (dec *decoder) alloc(size int) {
	l := len(dec.buf)
	if l == 0 {
		dec.buf = make([]byte, size)
	} else if size > l {
		if d := util.Abs(l - size); d > 0 {
			e := make([]byte, d)
			dec.buf = append(dec.buf, e...)
		}
	}
}

func (dec *decoder) readBytes(size int, new ...bool) ([]byte, error) {
	if size <= 0 {
		return nil, nil
	}

	dec.alloc(size)
	l, err := io.ReadFull(dec.r, dec.buf[:size])
	if err != nil {
		return nil, err
	}
	dec.read += l

	if len(new) > 0 && new[0] {
		data := make([]byte, size)
		copy(data, dec.buf)
		return data, nil
	}
	return dec.buf[:size], nil
}

func (dec *decoder) seek(size int) (err error) {
	_, err = dec.readBytes(size)
	return
}

func (dec *decoder) readPascalString() (string, int, error) {
	buf, err := dec.readBytes(1)
	if err != nil {
		return "", 0, err
	}

	size := int(buf[0])
	if size <= 0 {
		return "", 1, nil
	}

	buf, err = dec.readBytes(size)
	if err != nil {
		return "", size, err
	}

	return string(dec.buf[:size]), size, nil
}

func (dec *decoder) parseHeader() error {
	buf, err := dec.readBytes(headerLen)
	if err != nil {
		return err
	}

	// Signature
	read := headerLens[0]
	if !bytes.Equal(buf[:read], headerSig) {
		return ErrHeaderVersion
	}

	// Version
	dec.header.Version = int(util.ReadUint16(buf, read))
	read += headerLens[1]

	// Reserved: must be zero
	read += headerLens[2]

	// The number of channels in the image
	dec.header.Channels = int(util.ReadUint16(buf, read))
	read += headerLens[3]

	// The height of the image in pixels
	dec.header.Height = int(util.ReadUint32(buf, read))
	read += headerLens[4]

	// The width of the image in pixels
	dec.header.Width = int(util.ReadUint32(buf, read))
	read += headerLens[5]

	// Depth
	dec.header.Depth = int(util.ReadUint16(buf, read))
	read += headerLens[6]

	// The color mode of the file
	dec.header.ColorMode = ColorMode(util.ReadUint16(buf, read))
	read += headerLens[7]

	return nil
}

func (dec *decoder) parseColorModeData() (*ColorModeData, error) {
	buf, err := dec.readBytes(sectionLen)
	if err != nil {
		return nil, err
	}

	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, nil
	}

	if dec.header.ColorMode == ColorModeIndexed && size != 768 {
		return nil, ErrColorModeData
	}

	buf, err = dec.readBytes(size, true)
	if err != nil {
		return nil, err
	}

	return &ColorModeData{Data: buf}, nil
}

func (dec *decoder) parseImageResources() (blocks []*ImageResourceBlock, err error) {
	buf, err := dec.readBytes(sectionLen)
	if err != nil {
		return nil, err
	}

	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, nil
	}

	l := dec.read + size
	for dec.read < l {
		buf, err = dec.readBytes(sigLen)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(buf[:], imgResSig) {
			return nil, ErrImageResourceBlock
		}

		block := &ImageResourceBlock{}

		buf, err = dec.readBytes(uniqueIdentifierLen)
		if err != nil {
			return nil, err
		}
		block.ID = int(util.ReadUint16(buf, 0))

		str, l, err := dec.readPascalString()
		if err != nil {
			return nil, err
		}
		block.Name = str
		if l&1 != 0 {
			err = dec.seek(1)
			if err != nil {
				return nil, err
			}
		}

		buf, err = dec.readBytes(4)
		if err != nil {
			return nil, err
		}
		size = int(util.ReadUint32(buf, 0))
		buf, err = dec.readBytes(size, true)
		if err != nil {
			return nil, err
		}
		block.Data = buf

		if size&1 != 0 {
			err = dec.seek(1)
			if err != nil {
				return nil, err
			}
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (dec *decoder) parseLayer() ([]*Layer, error) {
	// TODO: adapt PSB
	buf, err := dec.readBytes(sectionLen)
	if err != nil {
		return nil, err
	}

	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, nil
	}

	layers, err := dec.parseLayerInfo()
	if err != nil {
		return nil, err
	}

	return layers, nil
}

func (dec *decoder) parseLayerInfo() ([]*Layer, error) {
	// TODO: PSB
	buf, err := dec.readBytes(4) // Length of the layers info section
	if err != nil {
		return nil, err
	}
	size := int(util.ReadUint32(buf, 0))
	size = size * 1

	buf, err = dec.readBytes(2)
	if err != nil {
		return nil, err
	}
	count := util.Abs(int(util.ReadInt16(buf, 0)))

	layers := make([]*Layer, count)
	for i := 0; i < count; i++ {
		layer, err := dec.parseLayerRecord()
		if err != nil {
			return nil, err
		}
		layer.Index = i
		layers[i] = layer
	}

	// Channel image data
	for i := range layers {
		layers[i].Image, err = dec.parseChannelImageData(layers[i])
		if err != nil {
			return nil, err
		}
	}

	return layers, nil
}

func (dec *decoder) parseLayerRecord() (*Layer, error) {
	buf, err := dec.readBytes(4*4 + 2)
	if err != nil {
		return nil, err
	}

	layer := newLayer()

	layer.setRect(
		int(util.ReadUint32(buf, 0)),
		int(util.ReadUint32(buf, 4)),
		int(util.ReadUint32(buf, 8)),
		int(util.ReadUint32(buf, 12)),
	)

	size := int(util.ReadUint16(buf, 16))
	layer.Channels = make([]*Channel, size)
	for i := range layer.Channels {
		channel := &Channel{}
		// TODO: PSB (2+8)
		buf, err := dec.readBytes(2 + 4)
		if err != nil {
			return nil, err
		}
		channel.ID = int(util.ReadInt16(buf, 0))
		channel.Length = int(util.ReadUint32(buf, 2))
		layer.Channels[i] = channel
	}

	buf, err = dec.readBytes(16)
	if !bytes.Equal(buf[0:4], layerSig) {
		return nil, ErrImageResourceBlock
	}
	layer.BlendModeKey = BlendModeKey(util.ReadString(buf, 4, 8))
	layer.Opacity = int(buf[8])
	layer.Clipping = Clipping(buf[9])
	layer.Flags = buf[10]
	layer.Filler = buf[11]

	size = int(util.ReadUint32(buf, 12))
	pos := size + dec.read

	// Mask
	layer.Mask, err = dec.parseMask()
	if err != nil {
		return nil, err
	}

	// Blending Ranges
	layer.BlendingRanges, err = dec.parseBlendingRanges()
	if err != nil {
		return nil, err
	}

	// Layer name (MBCS)
	var l int
	layer.LegacyName, l, err = dec.readPascalString()
	err = dec.seek((4 - ((1 + l) % 4)) % 4) // padding
	if err != nil {
		return nil, err
	}

	// Additional layer information
	for dec.read < pos {
		addInfo, err := dec.parseAdditionalLayerInfo()
		if err != nil {
			return nil, err
		}
		layer.setAdditionalInfo(addInfo)
	}

	return layer, nil
}

func (dec *decoder) parseMask() (*Mask, error) {
	buf, err := dec.readBytes(4)
	if err != nil {
		return nil, err
	}
	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, nil
	}

	buf, err = dec.readBytes(size)
	if err != nil {
		return nil, err
	}

	mask := newMask()

	mask.setRect(
		int(util.ReadUint32(buf, 0)),
		int(util.ReadUint32(buf, 4)),
		int(util.ReadUint32(buf, 8)),
		int(util.ReadUint32(buf, 12)),
	)

	// default color. 0 or 255
	if buf[16] != 0x00 && buf[16] != 0xff {
		return nil, errors.New("psd: invalid mask default color")
	}
	mask.DefaultColor = buf[16]
	mask.Flags = buf[17]

	if size == 20 {
		mask.Padding = int(util.ReadUint16(buf, 18))
	} else {
		mask.RealFlags = buf[18]
		if buf[19] != 0x00 && buf[19] != 0xff {
			return nil, errors.New("psd: invalid real user mask background")
		}
		mask.RealBackground = buf[19]
		mask.setRectEnclosingMask(
			int(util.ReadUint32(buf, 20)),
			int(util.ReadUint32(buf, 24)),
			int(util.ReadUint32(buf, 28)),
			int(util.ReadUint32(buf, 32)),
		)
	}

	return mask, nil
}

func (dec *decoder) parseBlendingRanges() (*BlendingRanges, error) {
	buf, err := dec.readBytes(4)
	if err != nil {
		return nil, err
	}
	size := int(util.ReadInt32(buf, 0))
	if size <= 0 {
		return nil, nil
	}

	buf, err = dec.readBytes(size)
	if err != nil {
		return nil, err
	}

	blendingRanges := newBlendingRanges()
	blendingRanges.CompositeGrayBlend = &BlendingRangesData{
		Source:      int(util.ReadUint32(buf[0:4], 0)),
		Destination: int(util.ReadUint32(buf[4:8], 0)),
	}
	l := size / 8
	for i := 1; i < l; i++ {
		data := &BlendingRangesData{
			Source:      int(util.ReadUint32(buf[i*4:(i+1)*4], 0)),
			Destination: int(util.ReadUint32(buf[(i+1)*4:(i+2)*4], 0)),
		}
		blendingRanges.addBlendingRangesData(data)
	}

	return blendingRanges, nil
}

func (dec *decoder) parseAdditionalLayerInfo() (*AdditionalInfo, error) {
	buf, err := dec.readBytes(4 * 3) // TODO: PSB
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(buf[0:4], layerSig) && !bytes.Equal(buf[0:4], additionalSig) {
		return nil, errors.New("psd: invalid addtional layer information")
	}

	addInfo := &AdditionalInfo{}

	addInfo.Key = util.ReadString(buf, 4, 8)

	size := int(util.ReadUint32(buf, 8))
	buf, err = dec.readBytes(size, true)
	if err != nil {
		return nil, err
	}
	addInfo.Data = buf

	return addInfo, nil
}

func (dec *decoder) parseChannelImageData(layer *Layer) (image.Image, error) {

	imgBuf := map[int][]byte{}
	for _, channel := range layer.Channels {
		buf, err := dec.readBytes(2)
		if err != nil {
			return nil, err
		}

		method := int(util.ReadUint16(buf, 0))

		if channel.Length == 2 {
			continue
		}

		var imgCh []byte
		switch method {
		case 0:
			imgCh, err = dec.parseRawImageData(layer)
		case 1:
			imgCh, err = dec.parseRLEImageData(layer)
		case 2, 3:

		default:

		}

		if err != nil {
			return nil, err
		}
		imgBuf[channel.ID] = imgCh
	}

	if len(imgBuf) <= 0 {
		return nil, nil
	}

	p := pixel.New(int(dec.header.ColorMode), dec.header.Depth, true)
	p.SetSource(
		layer.Rect.Min.Y,
		layer.Rect.Min.X,
		layer.Rect.Max.Y,
		layer.Rect.Max.X,
		imgBuf[0],
		imgBuf[1],
		imgBuf[2],
		imgBuf[-1],
	)

	return p, nil
}

func (dec *decoder) parseImageData() error {
	buf, err := dec.readBytes(compressionLen)
	if err != nil {
		return err
	}
	fmt.Println(buf)

	return nil
}

func (dec *decoder) parseRawImageData(layer *Layer) ([]byte, error) {
	size := layer.Rect.Dx() * layer.Rect.Dy()
	img, err := dec.readBytes(size, true)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (dec *decoder) parseRLEImageData(layer *Layer) ([]byte, error) {
	size := layer.Rect.Dy() * 4 / 2 // TODO: PSB 4 -> 8
	buf, err := dec.readBytes(size)
	if err != nil {
		return nil, err
	}
	lens := make([]int, layer.Rect.Dy())
	var total, n int
	for i := range lens {
		l := int(util.ReadUint16(buf, n))
		n += 2
		total += l
		lens[i] = l
	}
	buf, err = dec.readBytes(total, true)
	if err != nil {
		return nil, err
	}

	size = (layer.Rect.Dx()*dec.header.Depth + 7) >> 3 * layer.Rect.Dy()
	dest := make([]byte, size)
	decodePackBitsPerLine(dest, buf, lens)

	return dest, nil
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

func Decode(r io.Reader) (*PSD, error) {
	dec := &decoder{r: r, header: &Header{}}

	if err := dec.parseHeader(); err != nil {
		return nil, err
	}

	colorModeData, err := dec.parseColorModeData()
	if err != nil {
		return nil, err
	}

	blocks, err := dec.parseImageResources()
	if err != nil {
		return nil, err
	}

	layers, err := dec.parseLayer()

	// dec.parseImageData()

	psd := &PSD{
		Header:         dec.header,
		ColorModeData:  colorModeData,
		ImageResources: blocks,
		Layers:         layers,
	}
	return psd, nil
}
