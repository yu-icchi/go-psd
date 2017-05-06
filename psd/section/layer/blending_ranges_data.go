package layer

import (
	"github.com/yu-ichiko/go-psd/psd/util"
	"fmt"
)

func parseBlendingRangesData(buf []byte) (read int, err error) {
	size := int(util.ReadUint32(buf, 0))
	read += 4
	if size <= 0 {
		return
	}

	fmt.Println("------", buf[:size])
	read += size

	return
}
