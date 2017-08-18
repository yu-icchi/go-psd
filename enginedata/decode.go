package enginedata

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"regexp"
	"strconv"
	"unicode/utf16"
)

var (
	multiLineArrayStart = regexp.MustCompile(`^\/([a-zA-Z0-9]+) \[$`)
	property            = regexp.MustCompile(`^\/([a-zA-Z0-9]+)$`)
	propertyWithData    = regexp.MustCompile(`^\/([a-zA-Z0-9]+) (.*)$`)
	singleLineArray     = regexp.MustCompile(`^\[(.*)\]$`)
	boolean             = regexp.MustCompile(`^(true|false)$`)
	number              = regexp.MustCompile(`^(-?\d+)$`)
	numberWithDecimal   = regexp.MustCompile(`^(-?\d*)\.(\d+)$`)
	strRegexp           = regexp.MustCompile(`^\((.*)\)$`)
	crlf                = regexp.MustCompile(`\r\n|\r|\n|\f`)
	tab                 = []byte("\t")
	null                = []byte("")
	space               = []byte(" ")
	hashStart           = []byte("<<")
	hashEnd             = []byte(">>")
	multiLineArrayEnd   = []byte("]")
)

func Parser(buf []byte) (interface{}, error) {
	d := newDecoder()
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	for scanner.Scan() {
		b := bytes.Replace(scanner.Bytes(), tab, null, -1)
		_, err := d.parse(b)
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

func (d *decoder) parse(buf []byte) (interface{}, error) {
	buf = bytes.TrimSpace(buf)
	switch {
	case bytes.Equal(buf, hashStart):
		d.pushStack(object{})
	case bytes.Equal(buf, hashEnd):
		d.popStack()
	case multiLineArrayStart.Match(buf):
		data := multiLineArrayStart.FindSubmatch(buf)
		d.pushKeyStack(string(data[1]))
		d.pushStack(array{})
	case bytes.Equal(buf, multiLineArrayEnd):
		d.popStack()
	case property.Match(buf):
		data := property.FindSubmatch(buf)
		d.pushKeyStack(string(data[1]))
	case propertyWithData.Match(buf):
		data := propertyWithData.FindSubmatch(buf)
		v, err := d.parse(data[2])
		if err != nil {
			return nil, err
		}
		d.setCurrent(string(data[1]), v)
	case singleLineArray.Match(buf):
		data := singleLineArray.FindSubmatch(buf)
		arr := array{}
		for _, b := range bytes.Split(bytes.TrimSpace(data[1]), space) {
			v, err := d.parse(b)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case boolean.Match(buf):
		data := boolean.FindSubmatch(buf)
		return strconv.ParseBool(string(data[1]))
	case number.Match(buf):
		data := number.FindSubmatch(buf)
		return strconv.ParseInt(string(data[1]), 10, 32)
	case numberWithDecimal.Match(buf):
		data := numberWithDecimal.FindSubmatch(buf)
		d1 := string(data[1])
		d2 := string(data[2])
		if d1 == "" {
			d1 = "0"
		}
		return strconv.ParseFloat(d1+"."+d2, 64)
	case strRegexp.Match(buf):
		data := strRegexp.FindSubmatch(bytes.TrimSpace(crlf.ReplaceAll(buf, null)))
		str := decodeUTF16(data[1])
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
