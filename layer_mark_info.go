package psd

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf16"
)

func validBlendModeSignature(signature string) bool {
	return signature == "8BIM"
}

func getSize(is64 bool) int {
	if is64 {
		return 8
	}
	return 4
}

func readLayerAndMarkInfo(r io.Reader, psd *PSD) error {
	size := getSize(psd.IsPSB())
	buf := make([]byte, size)
	var read int

	l, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	read += l

	layerAndMarkInfoLen := int(readUint(buf))
	if layerAndMarkInfoLen == 0 {
		return nil
	}

	fmt.Println(layerAndMarkInfoLen)

	readLayerInfo(r, psd)

	return nil
}

func readLayerInfo(r io.Reader, psd *PSD) ([]Layer, int, error) {
	size := getSize(psd.IsPSB())
	buf := make([]byte, size)
	var read int

	l, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, 0, err
	}
	read += l

	layerInfoLen := int(readUint(buf))
	if layerInfoLen <= 0 {
		return nil, read, nil
	}

	fmt.Println("layerInfoLen:", layerInfoLen)

	buf = make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}

	numLayers := readUint16(buf, 0)
	fmt.Println("numLayers:", numLayers)

	readLayerRecords(r, 0, psd)

	return nil, read, nil
}

func readLayerRecords(r io.Reader, idx int, psd *PSD) (*Layer, int, error) {

	var read int
	buf := make([]byte, 4*4+2)
	l, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, read, err
	}
	read += l

	layer := Layer{
		Index: idx,
	}
	layer.Y = int(readUint32(buf, 0))
	layer.X = int(readUint32(buf, 4))
	layer.Height = int(readUint32(buf, 8))
	layer.Width = int(readUint32(buf, 12))

	numChannels := int(readUint16(buf, 4*4))
	fmt.Println("numChannels:", numChannels)

	channelInfoList := []ChannelInfo{}
	for i := 0; i < numChannels; i++ {
		size := getSize(psd.IsPSB())
		buf := make([]byte, 2+size)
		if l, err = io.ReadFull(r, buf); err != nil {
			return nil, read, err
		}
		read += l
		id := readUint16(buf, 0)
		var data uint64
		if psd.IsPSB() {
			data = readUint64(buf, 2)
		} else {
			data = uint64(readUint32(buf, 2))
		}
		ci := ChannelInfo{
			ID:   int(id),
			Data: data,
		}
		channelInfoList = append(channelInfoList, ci)
	}
	layer.ChannelInfo = channelInfoList

	buf = make([]byte, 12)
	if l, err = io.ReadFull(r, buf); err != nil {
		return nil, read, err
	}
	read += l

	blendModeSig := readString(buf, 0, 4)
	if !validBlendModeSignature(blendModeSig) {
		return nil, read, errors.New("psd: unexpected the blend mode signature")
	}
	layer.BlendMode = BlendMode{
		Signature: blendModeSig,
		Key:       readString(buf, 4, 8),
	}
	layer.Opacity = buf[8]
	layer.Clipping = buf[9]
	layer.Flags = buf[10]

	_, err = readLayerExtra(r)
	fmt.Println("readLayerExtra-err:", err)

	return &layer, 0, nil
}

func readLayerExtra(r io.Reader) (int, error) {
	var read int

	buf := make([]byte, 4)
	l, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	read += l

	extraDataLen := int(readUint32(buf, 0))
	if extraDataLen <= 0 {
		return read, nil
	}
	fmt.Println("extraDataLen:", extraDataLen)

	// Layer mask data
	if l, err = io.ReadFull(r, buf); err != nil {
		return read, err
	}
	read += l

	maskLen := int(readUint32(buf, 0))
	if maskLen > 0 {
		readMask(r)
	}
	fmt.Println("maskLen:", maskLen)

	// Layer blending ranges
	if l, err = io.ReadFull(r, buf); err != nil {
		return read, err
	}
	read += l

	blendingRangesLen := int(readUint32(buf, 0))
	fmt.Println("blendingRangesLen:", blendingRangesLen)
	if blendingRangesLen > 0 {
		buf := make([]byte, blendingRangesLen)
		if l, err = io.ReadFull(r, buf); err != nil {
			return read, err
		}
		read += l
		fmt.Println(buf)
	}

	pascalStr, read, err := readPascalStr(r)
	if err != nil {
		return read, err
	}
	fmt.Println(pascalStr, read)
	if l, err = adjustAlign4(r, read); err != nil {
		return read, err
	}
	read += l

	if read < extraDataLen+4 {
		l, err = addLayerInfo(r, extraDataLen+4-read)
		if err != nil {
			return read, err
		}
		read += l
	}

	return read, nil
}

func addLayerInfo(r io.Reader, len int) (int, error) {
	fmt.Println("=========== addLayerInfo", len)
	read := 0

	infoMap := map[string]AddLayerInfo{}
	for read < len {

		buf := make([]byte, 8)
		l, err := io.ReadFull(r, buf)
		if err != nil {
			return read, err
		}
		read += l
		fmt.Println(buf)

		sig := string(buf[:4])
		fmt.Println("sig", sig)
		if sig != "8BIM" && sig != "8B64" {
			buf = make([]byte, len - read)
			if l, err = io.ReadFull(r, buf); err != nil {
				return read, err
			}
			fmt.Println(buf)
			break
		}

		key := string(buf[4:8])

		buf = make([]byte, 4)
		if l, err = io.ReadFull(r, buf); err != nil {
			return read, err
		}
		read += l

		size := int(readUint32(buf, 0))
		fmt.Println("=== size:", size)
		buf = make([]byte, size)
		if l, err = io.ReadFull(r, buf); err != nil {
			return read, err
		}
		read += l

		ali := AddLayerInfo{
			Signature: sig,
			Key: key,
			Data: buf,
		}
		fmt.Println(ali, read)
		if key == "lyid" {
			fmt.Println(readUint32(buf, 0))
		}
		infoMap[key] = ali
	}

	fmt.Println(infoMap)
	return read, nil
}

func readUnicodeString(b []byte) string {
	size := readUint32(b, 0)
	if size <= 0 {
		return ""
	}
	buf := make([]uint16, size)
	for i := range buf {
		buf[i] = readUint16(b, 4+i<<1)
	}
	return string(utf16.Decode(buf))
}

func readMask(r io.Reader) {
	fmt.Println("mask")
}

func readPascalStr(r io.Reader) (string, int, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", 0, err
	}
	size := int(buf[0])
	fmt.Println(size)
	if size == 0 {
		return "", 1, nil
	}
	buf = make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", 1, err
	}

	return string(buf), size + 1, nil
}

func adjustAlign4(r io.Reader, l int) (int, error) {
	if gap := l & 3; gap > 0 {
		var b [4]byte
		return r.Read(b[:4-gap])
	}
	return 0, nil
}

func readAddLayerInfo(r io.Reader, infoLen int, psd *PSD) {
	read := 0
	buf := make([]byte, 8)
	for read < infoLen {
		buf[0] = 0
		for read < infoLen && buf[0] != '8' {

		}
	}
}
