package layer

import (
	"github.com/yu-ichiko/go-psd/util"
)

func parseMask(buf []byte) (read int, err error) {
	size := int(util.ReadUint32(buf, 0))
	read += 4
	if size <= 0 {
		return
	}

	// fmt.Println("== mask", size)
	mask := Mask{}

	mask.Top = int(util.ReadUint32(buf, read))
	read += 4
	mask.Left = int(util.ReadUint32(buf, read))
	read += 4
	mask.Bottom = int(util.ReadUint32(buf, read))
	read += 4
	mask.Right = int(util.ReadUint32(buf, read))
	read += 4

	mask.DefaultColor = int(util.ReadUint8(buf, read))
	read++

	/**
	 * bit 0 = position relative to layer
	 * bit 1 = layer mask disabled
	 * bit 2 = invert layer mask when blending (Obsolete)
	 * bit 3 = indicates that the user mask actually came from rendering other data
	 * bit 4 = indicates that the user and/or vector masks have parameters applied to them
	 */
	mask.Flags = buf[read]
	read++

	if size == 20 {
		// padding
		mask.Padding = int(util.ReadUint16(buf, read))
		read += 2
	} else {
		mask.RealFlags = int(util.ReadUint8(buf, read))
		read++

		mask.RealBackground = int(util.ReadUint8(buf, read))
		read++

		mask.RealTop = int(util.ReadUint32(buf, read))
		read += 4

		mask.RealLeft = int(util.ReadUint32(buf, read))
		read += 4

		mask.RealBottom = int(util.ReadUint32(buf, read))
		read += 4

		mask.RealRight = int(util.ReadUint32(buf, read))
		read += 4
	}

	// fmt.Printf("%+v\n", mask)
	return
}
