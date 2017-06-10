package layer

import (
	"errors"
	"fmt"
	"github.com/yu-ichiko/go-psd/section/header"
	"github.com/yu-ichiko/go-psd/util"
)

type AdditionalInfo struct {
	Key  string
	Data []byte
}

func parseAdditionalInfo(buf []byte, header *header.Header) (map[string]AdditionalInfo, error) {
	read := 0

	info := map[string]AdditionalInfo{}

	for read < len(buf) {
		sig := util.ReadString(buf, read, read+4)
		read += 4
		if sig != "8BIM" && sig != "8B64" {
			fmt.Println(sig)
			return info, errors.New("psd: invalid addtional layer information")
		}

		key := util.ReadString(buf, read, read+4)
		read += 4

		l := size(key, header.IsPSB())
		count := int(util.ReadUint(buf[read : read+l]))
		read += l

		data := buf[read : read+count]
		read += count

		add := AdditionalInfo{
			Key:  key,
			Data: data,
		}
		fmt.Println(add)
		info[key] = add
	}

	return info, nil
}

func size(key string, isPSB bool) int {
	if isPSB {
		switch key {
		case "LMsk",
			"Lr16",
			"Lr32",
			"Layr",
			"Mt16",
			"Mt32",
			"Mtrn",
			"Alph",
			"FMsk",
			"lnk2",
			"FEid",
			"FXid",
			"PxSD":
			return 8
		}
	}
	return 4
}
