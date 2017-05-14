package layer

import (
	"github.com/yu-ichiko/go-psd/util"
)

func parseBlendingRanges(buf []byte) (*BlendingRanges, int) {
	size := int(util.ReadUint32(buf, 0))
	read := 4
	if size <= 0 {
		return nil, read
	}

	blendingRanges := BlendingRanges{}

	blendingRanges.Black = int(util.ReadUint16(buf, read))
	read += 2

	blendingRanges.White = int(util.ReadUint16(buf, read))
	read += 2

	blendingRanges.DestRange = int(util.ReadUint32(buf, read))
	read += 4

	l := size - read + 4

	channels := make([]BlendingRangesChannel, l/8)
	for i := range channels {
		channel := BlendingRangesChannel{}

		channel.SourceRange = int(util.ReadUint32(buf, read))
		read += 4

		channel.DestinationRange = int(util.ReadUint32(buf, read))
		read += 4

		channels[i] = channel
	}
	blendingRanges.Channels = channels

	return &blendingRanges, read
}
