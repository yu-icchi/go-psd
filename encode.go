package psd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/yu-ichiko/go-psd/util"
	"io"
	"fmt"
	"image"
)

type encoder struct {
	w    *bufio.Writer
	buf  []byte
	read int
}

func (enc *encoder) alloc(size int) {
	l := len(enc.buf)
	if l == 0 {
		enc.buf = make([]byte, size)
	} else if size > l {
		if d := util.Abs(l - size); d > 0 {
			e := make([]byte, d)
			enc.buf = append(enc.buf, e...)
		}
	}
}

func (enc *encoder) writeBytes(buf []byte) error {
	n, err := enc.w.Write(buf)
	enc.read += n
	return err
}

func (enc *encoder) writeUint16(n int) error {
	enc.alloc(2)
	binary.BigEndian.PutUint16(enc.buf[:2], uint16(n))
	return enc.writeBytes(enc.buf[:2])
}

func (enc *encoder) writeUint32(n int) error {
	enc.alloc(4)
	binary.BigEndian.PutUint32(enc.buf[:4], uint32(n))
	return enc.writeBytes(enc.buf[:4])
}

func (enc *encoder) seek(size int) (err error) {
	enc.alloc(size)
	for i := range enc.buf[:size] {
		enc.buf[i] = 0
	}
	err = enc.writeBytes(enc.buf[:size])
	return
}

func (enc *encoder) composeHeader(h *Header) error {
	// Signature
	if err := enc.writeBytes(headerSig); err != nil {
		return err
	}
	// Version
	if err := enc.writeUint16(h.Version); err != nil {
		return err
	}
	// Reserved
	if err := enc.seek(headerLens[2]); err != nil {
		return err
	}
	// Channels
	if err := enc.writeUint16(h.Channels); err != nil {
		return err
	}
	// Height
	if err := enc.writeUint32(h.Height); err != nil {
		return err
	}
	// Width
	if err := enc.writeUint32(h.Width); err != nil {
		return err
	}
	// Depth
	if err := enc.writeUint16(h.Depth); err != nil {
		return err
	}
	// Color Mode
	if err := enc.writeUint16(int(h.ColorMode)); err != nil {
		return err
	}
	return nil
}

func (enc *encoder) composeColorModeData(colorMode *ColorModeData) error {
	var size int
	if colorMode != nil {
		size = len(colorMode.Data)
	}
	if err := enc.writeUint32(size); err != nil {
		return err
	}
	if size > 0 {
		if err := enc.writeBytes(colorMode.Data); err != nil {
			return err
		}
	}
	return nil
}

func (enc *encoder) composeImageResources(blocks []*ImageResourceBlock) error {
	buf := &bytes.Buffer{}
	var size int
	for _, block := range blocks {
		if _, err := buf.Write(imgResSig); err != nil {
			return err
		}
		if _, err := buf.Write(util.ByteUint16(block.ID)); err != nil {
			return err
		}
		nameBuf, l, err := util.BytePascalString(block.Name)
		if err != nil {
			return err
		}
		if _, err := buf.Write(nameBuf); err != nil {
			return err
		}
		if l&1 != 0 {
			if _, err := buf.Write([]byte{0}); err != nil {
				return err
			}
		}
		size = len(block.Data)
		if _, err := buf.Write(util.ByteUint32(size)); err != nil {
			return err
		}
		if _, err := buf.Write(block.Data); err != nil {
			return err
		}
		if size&1 != 0 {
			if _, err := buf.Write([]byte{0}); err != nil {
				return err
			}
		}
	}
	if err := enc.writeUint32(buf.Len()); err != nil {
		return err
	}
	if err := enc.writeBytes(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (enc *encoder) composeLayerAndMaskInfo(layers []*Layer, mask *GlobalLayerMask, info []*AdditionalInfo) error {
	enc.composeLayerInfo(layers)
	return nil
}

func (enc *encoder) composeLayerInfo(layers []*Layer) error {
	layerRecordBuf := &bytes.Buffer{}
	channelImageBuf := &bytes.Buffer{}

	for _, layer := range layers {
		// Layer Record
		buf, err := enc.composeLayerRecord(layer)
		if err != nil {
			return err
		}
		layerRecordBuf.Write(buf.Bytes())

		// fmt.Println("==========")
		// Channel image data
		// enc.composeChannelImageData(layer.Image)
	}

	fmt.Println("layer record ===>", layerRecordBuf.Len())
	fmt.Println("channel image ===>", channelImageBuf.Len())

	return nil
}

func (enc *encoder) composeLayerRecord(layer *Layer) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	// Rect
	// Top
	if _, err := buf.Write(util.ByteUint32(layer.Rect.Min.Y)); err != nil {
		return nil, err
	}
	// Left
	if _, err := buf.Write(util.ByteUint32(layer.Rect.Min.X)); err != nil {
		return nil, err
	}
	// Bottom
	if _, err := buf.Write(util.ByteUint32(layer.Rect.Max.Y)); err != nil {
		return nil, err
	}
	// Left
	if _, err := buf.Write(util.ByteUint32(layer.Rect.Max.X)); err != nil {
		return nil, err
	}
	// Number of channels
	if _, err := buf.Write(util.ByteUint16(len(layer.Channels))); err != nil {
		return nil, err
	}
	// Channel information
	for _, channel := range layer.Channels {
		if _, err := buf.Write(util.ByteUint16(channel.ID)); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(channel.Length)); err != nil {
			return nil, err
		}
	}
	// Blend mode signature
	if _, err := buf.Write(layerSig); err != nil {
		return nil, err
	}
	// Blend mode key
	if _, err := buf.Write([]byte(string(layer.BlendModeKey))); err != nil {
		return nil, err
	}
	// Opacity, Clipping, Flags, Filler
	if _, err := buf.Write([]byte{byte(layer.Opacity), byte(layer.Clipping), layer.Flags, layer.Filler}); err != nil {
		return nil, err
	}

	// Mask
	maskBuf, err := enc.composeMask(layer.Mask)
	if err != nil {
		return nil, err
	}

	// Blending Ranges
	blendRangesBuf, err := enc.composeBlendingRanges(layer.BlendingRanges)
	if err != nil {
		return nil, err
	}

	// Layer name (MBCS)
	nameBuf, l, err := util.BytePascalString(layer.LegacyName)
	if err != nil {
		return nil, err
	}
	pl := (4 - ((1 + l) % 4)) % 4

	// Additional layer information
	addInfosBuf, err := enc.composeAdditionalLayerInfo(layer.AdditionalInfos)
	if err != nil {
		return nil, err
	}

	// extra data length
	size := maskBuf.Len() + 4
	size += blendRangesBuf.Len() + 4
	size += l + pl + 1
	size += addInfosBuf.Len()
	if _, err := buf.Write(util.ByteUint32(size)); err != nil {
		return nil, err
	}

	// Mask
	size = maskBuf.Len()
	if _, err := buf.Write(util.ByteUint32(size)); err != nil {
		return nil, err
	}
	if size > 0 {
		if _, err := buf.Write(buf.Bytes()); err != nil {
			return nil, err
		}
	}

	// Blending Ranges
	size = blendRangesBuf.Len()
	if _, err := buf.Write(util.ByteUint32(size)); err != nil {
		return nil, err
	}
	if size > 0 {
		if _, err := buf.Write(buf.Bytes()); err != nil {
			return nil, err
		}
	}

	// Layer name (MBCS) set
	if _, err := buf.Write(nameBuf); err != nil {
		return nil, err
	}
	// name padding
	if _, err := buf.Write(make([]byte, pl)); err != nil {
		return nil, err
	}

	// Additional layer information
	if addInfosBuf != nil {
		if _, err := buf.Write(addInfosBuf.Bytes()); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func (enc *encoder) composeMask(mask *Mask) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if mask == nil {
		return buf, nil
	}
	if _, err := buf.Write(util.ByteUint32(mask.Rect.Min.Y)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(util.ByteUint32(mask.Rect.Min.X)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(util.ByteUint32(mask.Rect.Max.Y)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(util.ByteUint32(mask.Rect.Max.X)); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte{mask.DefaultColor, mask.Flags}); err != nil {
		return nil, err
	}
	if len(mask.Padding) > 0 {
		if _, err := buf.Write(mask.Padding); err != nil {
			return nil, err
		}
	} else {
		if _, err := buf.Write([]byte{*mask.RealFlags, *mask.RealBackground}); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(mask.RectEnclosingMask.Min.Y)); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(mask.RectEnclosingMask.Min.X)); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(mask.RectEnclosingMask.Max.Y)); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(mask.RectEnclosingMask.Max.X)); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (enc *encoder) composeBlendingRanges(blendingRanges *BlendingRanges) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if blendingRanges == nil {
		return buf, nil
	}
	if _, err := buf.Write(util.ByteUint32(blendingRanges.CompositeGrayBlend.Source)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(util.ByteUint32(blendingRanges.CompositeGrayBlend.Destination)); err != nil {
		return nil, err
	}
	for _, channel := range blendingRanges.Channels {
		if _, err := buf.Write(util.ByteUint32(channel.Source)); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(channel.Destination)); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (enc *encoder) composeAdditionalLayerInfo(infos []*AdditionalInfo) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	for _, info := range infos {
		if _, err := buf.Write(additionalSig); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(info.Key); err != nil {
			return nil, err
		}
		if _, err := buf.Write(util.ByteUint32(len(info.Data))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(info.Data); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (enc *encoder) composeChannelImageData(img image.Image) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			//r, g, b, a := img.At(x, y).RGBA()
		}
	}
}

func Encode(w io.Writer, psd *PSD) error {
	bw := bufio.NewWriter(w)
	enc := &encoder{w: bw}

	var err error
	err = enc.composeHeader(psd.Header)
	if err != nil {
		return err
	}

	err = enc.composeColorModeData(psd.ColorModeData)
	if err != nil {
		return err
	}

	err = enc.composeImageResources(psd.ImageResources)
	if err != nil {
		return err
	}

	err = enc.composeLayerAndMaskInfo(psd.Layers, psd.GlobalLayerMask, psd.AdditionalInfos)
	if err != nil {
		return err
	}

	if err := enc.w.Flush(); err != nil {
		return err
	}
	return nil
}
