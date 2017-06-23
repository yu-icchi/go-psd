package psd

import (
	"bufio"
	"encoding/binary"
	"github.com/yu-ichiko/go-psd/util"
	"io"
	"fmt"
	"bytes"
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
	size := 0
	if colorMode != nil {
		size = len(colorMode.Data)
	}
	if err := enc.writeUint32(size); err != nil {
		return err
	}

	if size > 0 {
		// TODO...
	}

	return nil
}

func (enc *encoder) composeImageResources(blocks []*ImageResourceBlock) error {
	buf := &bytes.Buffer{}

	for _, block := range blocks {
		buf.Write(imgResSig)
		buf.Write(util.ByteUint16(block.ID))
		buf.Write(util.BytePascalString(block.Name))
		l := len(block.Name)
		if l == 0 {
			l = 1
		}
		if l&1 != 0 {
			buf.Write([]byte{0})
		}
	}
	fmt.Println(buf.Len())
	fmt.Println(buf.Bytes())

	return nil
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

	if err := enc.w.Flush(); err != nil {
		return err
	}
	return nil
}
