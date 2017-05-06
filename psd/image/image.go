package image

import (
	"io"

	"github.com/yu-ichiko/go-psd/psd/util"

	"fmt"
)

func Parse(r io.Reader) (read int, err error) {
	var l int
	buf := make([]byte, 2)
	if l, err = io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	read += l

	method := int(util.ReadUint16(buf, 0))
	fmt.Println(method)

	return
}
