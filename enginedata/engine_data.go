package enginedata

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

var (
	hashStart           = regexp.MustCompile(`^<<$`)
	hashEnd             = regexp.MustCompile(`^>>$`)
	multiLineArrayStart = regexp.MustCompile(`^\/([a-zA-Z0-9]+) \[$`)
	multiLineArrayEnd   = regexp.MustCompile(`^\]$`)
	property            = regexp.MustCompile(`^\/([a-zA-Z0-9]+)$`)
	propertyWithData    = regexp.MustCompile(`^\/([a-zA-Z0-9]+) (.*)$`)
	singleLineArray     = regexp.MustCompile(`^\[(.*)\]$`)
	boolean             = regexp.MustCompile(`^(true|false)$`)
	number              = regexp.MustCompile(`^(-?\d+)$`)
	numberWithDecimal   = regexp.MustCompile(`^(-?\d*)\.(\d+)$`)
	strRegexp           = regexp.MustCompile(`^\((.*)\)$`)

	crlf = regexp.MustCompile(`\r\n|\r|\n|\f`)
)

func Decode(buf []byte) (interface{}, error) {
	d := newDecoder()
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	for scanner.Scan() {
		t := strings.Replace(scanner.Text(), "\t", "", -1)
		_, err := d.parse(t)
		if err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return d.current, nil
}

func newDecoder() *decoder {
	return &decoder{
		current:  nil,
		keyStack: []string{},
		stack:    []interface{}{},
	}
}

type object map[string]interface{}
type array []interface{}

type decoder struct {
	current  interface{}
	keyStack []string
	stack    []interface{}
}

func (d *decoder) pushStack(node interface{}) {
	d.stack = append(d.stack, d.current)
	d.current = node
}

func (d *decoder) popStack() {
	// pop
	l := len(d.stack) - 1
	node := d.stack[l]
	d.stack = d.stack[:l]

	switch node.(type) {
	case array:
		arr := node.(array)
		arr = append(arr, d.current)
		d.current = arr
	case object:
		k := d.popKeyStack()
		obj := node.(object)
		obj[k] = d.current
		d.current = obj
	}
}

func (d *decoder) pushKeyStack(key string) {
	d.keyStack = append(d.keyStack, key)
}

func (d *decoder) popKeyStack() string {
	key := d.keyStack[len(d.keyStack)-1]
	d.keyStack = d.keyStack[:len(d.keyStack)-1]
	return key
}

func (d *decoder) setCurrent(key string, value interface{}) {
	obj := d.current.(object)
	obj[key] = value
	d.current = obj
}

func (d *decoder) parse(str string) (interface{}, error) {
	str = strings.TrimSpace(str)
	switch {
	case hashStart.MatchString(str):
		d.pushStack(object{})
	case hashEnd.MatchString(str):
		d.popStack()
	case multiLineArrayStart.MatchString(str):
		data := multiLineArrayStart.FindStringSubmatch(str)
		d.pushKeyStack(data[1])
		d.pushStack(array{})
	case multiLineArrayEnd.MatchString(str):
		d.popStack()
	case property.MatchString(str):
		data := property.FindStringSubmatch(str)
		d.pushKeyStack(data[1])
	case propertyWithData.MatchString(str):
		data := propertyWithData.FindStringSubmatch(str)
		v, err := d.parse(data[2])
		if err != nil {
			return nil, err
		}
		d.setCurrent(data[1], v)
	case singleLineArray.MatchString(str):
		data := singleLineArray.FindStringSubmatch(str)
		arr := array{}
		for _, ss := range strings.Split(strings.TrimSpace(data[1]), " ") {
			v, err := d.parse(ss)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case boolean.MatchString(str):
		data := boolean.FindStringSubmatch(str)
		return strconv.ParseBool(data[1])
	case number.MatchString(str):
		data := number.FindStringSubmatch(str)
		return strconv.ParseInt(data[1], 10, 32)
	case numberWithDecimal.MatchString(str):
		data := numberWithDecimal.FindStringSubmatch(str)
		if data[1] == "" {
			data[1] = "0"
		}
		return strconv.ParseFloat(data[1]+"."+data[2], 64)
	case strRegexp.MatchString(str):
		data := strRegexp.FindStringSubmatch(strings.TrimSpace(crlf.ReplaceAllString(str, "")))
		str = decodeUTF16([]byte(data[1]))
		return str, nil
	}
	return nil, nil
}

func decodeUTF16(buf []byte) string {
	data := make([]uint16, 0, len(buf)/2)
	for i := 0; i < len(buf); i += 2 {
		data = append(data, binary.BigEndian.Uint16(buf[i:i+2]))
	}
	return string(utf16.Decode(data))
}
