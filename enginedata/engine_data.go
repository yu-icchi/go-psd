package enginedata

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	hashStart           = regexp.MustCompile(`^<<$`)
	hashEnd             = regexp.MustCompile(`^>>$`)
	multiLineArrayStart = regexp.MustCompile(`^\/(\w+) \[$`)
	multiLineArrayEnd   = regexp.MustCompile(`^\]$`)
	property            = regexp.MustCompile(`^\/([a-zA-Z0-9]+)$`)
	propertyWithData    = regexp.MustCompile(`^\/([a-zA-Z0-9]+)\s((.|\r)*)$`)
	singleLineArray     = regexp.MustCompile(`^\[(.*)\]$`)
	boolean             = regexp.MustCompile(`^(true|false)$`)
	number              = regexp.MustCompile(`^-?\d+$`)
	numberWithDecimal   = regexp.MustCompile(`^(-?\d*)\.(\d+)$`)
	str                 = regexp.MustCompile(`^\(((.|\r)*)\)$`)
)

func Analyse(buf []byte) {
	text := strings.Split(string(buf), "\n")
	for _, t := range text {
		s := strings.Replace(t, "\t", "", -1)
		parser(s)
	}
}

func parser(s string) {
	switch {
	case hashStart.MatchString(s):
		fmt.Println("------>", "hash start", s)
	case hashEnd.MatchString(s):
		fmt.Println("------>", "hash end", s)
	case multiLineArrayStart.MatchString(s):
		fmt.Println("------>", "multi line array start", s)
	case multiLineArrayEnd.MatchString(s):
		fmt.Println("------>", "multi line array end", s)
	case property.MatchString(s):
		fmt.Println("------>", "property", s)
	case propertyWithData.MatchString(s):
		fmt.Println("------>", "propertyWithData", s)
		idx := strings.Index(s, " ")
		parser(s[idx+1:])
	case singleLineArray.MatchString(s):
		fmt.Println("------>", "singleLineArray", s)
	case boolean.MatchString(s):
		fmt.Println("------>", "boolean", s)
	case number.MatchString(s):
		fmt.Println("------>", "number", s)
	case numberWithDecimal.MatchString(s):
		fmt.Println("------>", "numberWithDecimal", s)
	case str.MatchString(s):
		fmt.Println("------>", "string", s)
	}
}
