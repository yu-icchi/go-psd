package psd

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"

	"encoding/json"
	"github.com/yu-ichiko/go-psd/enginedata"
	"github.com/yu-ichiko/go-psd/util"
)

type decoder struct {
	r    io.Reader
	buf  []byte
	read int

	header *Header
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

func (dec *decoder) readPackBits(l int) ([]byte, error) {

	limit := dec.read + l
	data := []byte{}

	// decode
	for dec.read < limit {
		buf, err := dec.readBytes(1)
		if err != nil {
			return nil, err
		}
		run := int(int8(buf[0]))
		if run < 0 {
			run = 1 - run
			buf, err = dec.readBytes(1, true)
			if err != nil {
				return nil, err
			}
			for run > 0 {
				data = append(data, buf...)
				run -= 1
			}
		} else {
			run = 1 + run
			for run > 0 {
				buf, err = dec.readBytes(1, true)
				if err != nil {
					return nil, err
				}
				data = append(data, buf...)
				run -= 1
			}
		}
	}

	return data, nil
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

		buf, err = dec.readBytes(actualLen)
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

func (dec *decoder) parseLayerAndMaskInfo() ([]*Layer, *GlobalLayerMask, []*AdditionalInfo, error) {
	s := dec.read

	// TODO: adapt PSB
	buf, err := dec.readBytes(sectionLen)
	if err != nil {
		return nil, nil, nil, err
	}

	size := int(util.ReadUint32(buf, 0))
	if size <= 0 {
		return nil, nil, nil, nil
	}
	pos := dec.read + size

	// Layer Info
	layers, err := dec.parseLayerInfo()
	if err != nil {
		return nil, nil, nil, err
	}

	// padding
	if padding := (dec.read - s + 4 - size) & 3; padding > 0 {
		err = dec.seek(4 - padding)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// Global Layer Info
	globalMask, err := dec.parseGlobalLayerMask()
	if err != nil {
		return nil, nil, nil, err
	}

	// Additional Layer Info
	addInfos := []*AdditionalInfo{}
	for dec.read < pos {
		addInfo, err := dec.parseAdditionalLayerInfo()
		if err != nil {
			return nil, nil, nil, err
		}
		addInfos = append(addInfos, addInfo)
	}

	return layers, globalMask, addInfos, nil
}

func (dec *decoder) parseLayerInfo() ([]*Layer, error) {
	// TODO: PSB
	buf, err := dec.readBytes(4) // Length of the layers info section
	if err != nil {
		return nil, err
	}
	size := int(util.ReadUint32(buf, 0))
	size = size * 1 // TODO: PSB is *2
	fmt.Println("size ===>", size)

	buf, err = dec.readBytes(2)
	if err != nil {
		return nil, err
	}
	count := util.Abs(int(util.ReadInt16(buf, 0)))
	fmt.Println("count ===>", buf, count, int(util.ReadInt16(buf, 0)))

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

	// Flags:
	// bit 0 = transparency protected;
	// bit 1 = visible;
	// bit 2 = obsolete;
	// bit 3 = 1 for Photoshop 5.0 and later, tells if bit 4 has useful information;
	// bit 4 = pixel data irrelevant to appearance of document
	layer.TransparencyProtected = (layer.Flags & (1 << 0)) == 0
	layer.Visible = (layer.Flags & (1 << 1)) == 0
	layer.Obsolete = (layer.Flags & (1 << 2)) == 0
	layer.IrrelevantPixelData = (layer.Flags & (1 << 4)) == 0

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

	buf, err = dec.readBytes(size, true)
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
		mask.Padding = buf[18:20]
	} else {
		mask.RealFlags = &buf[18]
		if buf[19] != 0x00 && buf[19] != 0xff {
			return nil, errors.New("psd: invalid real user mask background")
		}
		mask.RealBackground = &buf[19]
		mask.setRectEnclosingMask(
			int(util.ReadUint32(buf, 20)),
			int(util.ReadUint32(buf, 24)),
			int(util.ReadUint32(buf, 28)),
			int(util.ReadUint32(buf, 32)),
		)
	}

	return mask, nil
}

func (dec *decoder) parseGlobalLayerMask() (*GlobalLayerMask, error) {
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

	colorComponents := make([]byte, 8)
	for i, d := range buf[2:10] {
		colorComponents[i] = d
	}

	mask := GlobalLayerMask{}
	mask.OverlayColor = int(util.ReadUint16(buf, 0))
	mask.ColorComponents = colorComponents
	mask.Opacity = int(util.ReadUint16(buf, 10))
	mask.Kind = int(buf[12])
	mask.Fillers = len(buf[13:])

	return &mask, nil
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
		return nil, errors.New("psd: invalid additional layer information")
	}

	addInfo := &AdditionalInfo{}
	addInfo.Key = util.ReadString(buf, 4, 8)

	size := int(util.ReadUint32(buf, 8))

	// padding
	switch addInfo.Key {
	case "Txt2":
		if size%2 == 1 {
			size += 1
		}
		size += 2
	case "LMsk":
		size += 2
	}

	buf, err = dec.readBytes(size, true)
	if err != nil {
		return nil, err
	}
	addInfo.Data = buf
	if addInfo.Key == "TySh" {
		dec.parseTypeToolObjectSetting(buf)
	}

	return addInfo, nil
}

func (dec *decoder) parseTypeToolObjectSetting(buf []byte) {
	read := 0
	fmt.Println("==== type tool object setting start ====")
	fmt.Println(buf)
	version := util.ReadInt16(buf, read)
	read += 2
	fmt.Println("version:", version)
	xx := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("xx:", xx)
	xy := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("xy:", xy)
	yx := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("yx:", yx)
	yy := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("yy:", yy)
	tx := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("tx:", tx)
	ty := util.ReadUint64(buf, read)
	read += 8
	fmt.Println("ty:", ty)
	textVersion := util.ReadUint16(buf, read)
	read += 2
	fmt.Println("text version:", textVersion)
	descriptionVersion := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("description version", descriptionVersion)

	classID, l := util.UnicodeString(buf[read:])
	fmt.Println("unicode classID:", classID, l)
	read += l
	str, l := util.ReadClassID(buf[read:])
	fmt.Println("classID:", str, l)
	read += l
	num := int(util.ReadUint32(buf, read))
	fmt.Println("descriptor num ====>", num)
	read += 4
	for i := 0; i < num; i++ {
		id, l := util.ReadClassID(buf[read:])
		fmt.Println("Text data id =====>", id, l)
		read += l
		osTypeKey := util.ReadString(buf, read, read+4)
		read += 4
		fmt.Println("osTypeKey:", osTypeKey)
		switch osTypeKey {
		case "TEXT":
			value, l := util.UnicodeString(buf[read:])
			fmt.Println("TEXT -------->", value)
			read += l
		case "enum":
			id, l := util.ReadClassID(buf[read:])
			fmt.Println("enum id ------->", id, l)
			read += l
			id, l = util.ReadClassID(buf[read:])
			fmt.Println("enum value ------->", id, l)
			read += l
		case "long":
			num := int(util.ReadUint32(buf, read))
			fmt.Println("long ------->", num)
			read += 4
		case "tdta":
			size := int(util.ReadUint32(buf, read))
			fmt.Println("tdta size ------->", size)
			read += 4
			data, _ := enginedata.Decode(buf[read : read+size])
			jb, _ := json.Marshal(data)
			fmt.Println("tdta ------->", string(jb))
			read += size
		case "bool":
			fmt.Println("bool ------->", buf[read])
			read += 1
		}
	}

	warpVersion := util.ReadInt16(buf, read)
	read += 2
	fmt.Println("warpVersion:", warpVersion)
	descriptorVersion := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("descriptorVersion:", descriptorVersion)

	classID, l = util.UnicodeString(buf[read:])
	fmt.Println("unicode classID:", classID, l)
	read += l
	str, l = util.ReadClassID(buf[read:])
	fmt.Println("classID:", str, l)
	read += l
	num = int(util.ReadUint32(buf, read))
	fmt.Println("descriptor num ====>", num)
	read += 4
	for i := 0; i < num; i++ {
		id, l := util.ReadClassID(buf[read:])
		fmt.Println("id =====>", id, l)
		read += l
		osTypeKey := util.ReadString(buf, read, read+4)
		read += 4
		fmt.Println("osTypeKey:", osTypeKey)
		switch osTypeKey {
		case "TEXT":
			value, l := util.UnicodeString(buf[read:])
			fmt.Println("TEXT -------->", value)
			read += l
		case "enum":
			id, l := util.ReadClassID(buf[read:])
			fmt.Println("enum id ------->", id, l)
			read += l
			id, l = util.ReadClassID(buf[read:])
			fmt.Println("enum value ------->", id, l)
			read += l
		case "long":
			num := int(util.ReadUint32(buf, read))
			fmt.Println("long ------->", num)
			read += 4
		case "tdta":
			size := int(util.ReadUint32(buf, read))
			fmt.Println("tdta size ------->", size)
			read += 4
			fmt.Println("tdta ------->", buf[read:size])
			read += size
		case "bool":
			fmt.Println("bool ------->", buf[read])
			read += 1
		case "doub":
			num := util.ReadUint64(buf, read)
			fmt.Println("doub ------->", num)
			read += 8
		}
	}

	left := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("left:", left)
	top := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("top:", top)
	right := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("right:", right)
	bottom := util.ReadUint32(buf, read)
	read += 4
	fmt.Println("bottom:", bottom)

	fmt.Println("==== type tool object setting end ====", len(buf)-read)
}

func (dec *decoder) parseChannelImageData(layer *Layer) (Image, error) {
	var method int
	img := map[int][]byte{}
	for _, channel := range layer.Channels {
		buf, err := dec.readBytes(compressionLen)
		if err != nil {
			return nil, err
		}

		method = int(util.ReadUint16(buf, 0))

		if channel.Length == 2 {
			continue
		}

		var rect image.Rectangle
		switch channel.ID {
		case -3:
			rect = *layer.Mask.RectEnclosingMask
		case -2:
			rect = layer.Mask.Rect
		default:
			rect = layer.Rect
		}

		var imgCh []byte
		switch method {
		case imgRAW:
			// Raw Image
			imgCh, err = dec.parseChannelImageRaw(rect)
		case imgRLE:
			// RLE
			imgCh, err = dec.parseChannelImageRLE(rect)
		default:
			return nil, fmt.Errorf("psd: unknown compression method=%d", method)
		}

		if err != nil {
			return nil, err
		}
		if channel.ID >= -1 {
			img[channel.ID] = imgCh
		}
	}

	if len(img) <= 0 {
		return nil, nil
	}

	hasAlpha := dec.header.ColorMode.Channels() < len(img)
	p, err := newImage(dec.header.ColorMode, dec.header.Depth, method, hasAlpha)
	if err != nil {
		return nil, err
	}
	p.Source(layer.Rect, img[0], img[1], img[2], img[-1])

	return p, nil
}

func (dec *decoder) parseChannelImageRaw(rect image.Rectangle) ([]byte, error) {
	size := rect.Dx() * rect.Dy()
	img, err := dec.readBytes(size, true)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (dec *decoder) parseChannelImageRLE(rect image.Rectangle) ([]byte, error) {
	size := rect.Dy() * 4 / 2 // TODO: PSB 4 -> 8
	buf, err := dec.readBytes(size)
	if err != nil {
		return nil, err
	}
	lens := make([]int, rect.Dy())
	var total, n int
	for i := range lens {
		l := int(util.ReadUint16(buf, n))
		n += 2
		total += l
		lens[i] = l
	}
	buf, err = dec.readBytes(total)
	if err != nil {
		return nil, err
	}

	size = (rect.Dx()*dec.header.Depth + 7) >> 3 * rect.Dy()
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

func (dec *decoder) parseImageData() (image.Image, error) {
	buf, err := dec.readBytes(compressionLen)
	if err != nil {
		return nil, err
	}

	method := int(util.ReadUint16(buf, 0))

	switch method {
	case imgRAW:
		return dec.parseImageRAW()
	case imgRLE:
		return dec.parseImageRLE()
	default:
		return nil, fmt.Errorf("psd: unknown compression method=%d", method)
	}
	return nil, nil
}

func (dec *decoder) parseImageRAW() (Image, error) {
	size := dec.header.Width * dec.header.Height * (dec.header.Depth / 8)
	img := make([][]byte, dec.header.Channels)
	var err error

	for i := 0; i < dec.header.Channels; i++ {
		img[i], err = dec.readBytes(size, true)
		if err != nil {
			return nil, err
		}
	}

	p, err := newImage(dec.header.ColorMode, dec.header.Depth, imgRAW, false)
	if err != nil {
		return nil, err
	}
	p.Source(dec.header.Rect(), img...)

	return p, nil
}

func (dec *decoder) parseImageRLE() (Image, error) {

	lineLen := make([]int, dec.header.Height*dec.header.Channels)
	for i := range lineLen {
		buf, err := dec.readBytes(2)
		if err != nil {
			return nil, err
		}
		lineLen[i] = int(util.ReadUint16(buf, 0))
	}

	img := make([][]byte, dec.header.Channels)
	d := dec.header.Depth / 8
	for i := 0; i < dec.header.Channels; i++ {
		lines := make([]byte, 0, dec.header.Height)
		size := 0
		for j := 0; j < dec.header.Height; j++ {
			n := lineLen[i*dec.header.Height+j] * d
			line, err := dec.readPackBits(n)
			if err != nil {
				return nil, err
			}
			lines = append(lines, line...)
			size += len(line)
		}
		img[i] = lines
	}

	p, err := newImage(dec.header.ColorMode, dec.header.Depth, imgRLE, false)
	if err != nil {
		return nil, err
	}
	p.Source(dec.header.Rect(), img...)

	return p, nil
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

	layers, globalMask, addInfos, err := dec.parseLayerAndMaskInfo()
	if err != nil {
		return nil, err
	}

	img, err := dec.parseImageData()
	if err != nil {
		return nil, err
	}

	psd := &PSD{
		Header:          dec.header,
		ColorModeData:   colorModeData,
		ImageResources:  blocks,
		Layers:          layers,
		GlobalLayerMask: globalMask,
		AdditionalInfos: addInfos,
		Image:           img,
	}
	return psd, nil
}
