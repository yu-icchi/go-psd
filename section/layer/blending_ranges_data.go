package layer

import (
	"github.com/yu-ichiko/go-psd/util"
)

func parseBlendingRangesData(buf []byte) (read int, err error) {
	size := int(util.ReadUint32(buf, 0))
	read += 4
	if size <= 0 {
		return
	}

	//black := util.ReadUint16(buf, read)
	//read += 2
	// fmt.Println("black", black)

	//white := util.ReadUint16(buf, read)
	//read += 2
	// fmt.Println("white", white)

	//destRange := util.ReadUint32(buf, read)
	//read += 4
	// fmt.Println("destRange", destRange)

	// fmt.Println("----->", size - read)

	return
}
