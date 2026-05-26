package cmdutil

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var utf8BOM = []byte{0xef, 0xbb, 0xbf}
var utf16LEBOM = []byte{0xff, 0xfe}
var utf16BEBOM = []byte{0xfe, 0xff}

// ReadTextFile reads a user-provided text file and strips a UTF-8 BOM when present.
func ReadTextFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return DecodeUserText(content), nil
}

// ReadText reads user-provided text from a stream.
func ReadText(r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return DecodeUserText(content), nil
}

// DecodeUserText decodes text accepted from user-facing files or stdin.
func DecodeUserText(content []byte) string {
	content = bytes.TrimPrefix(content, utf8BOM)
	if bytes.HasPrefix(content, utf16LEBOM) {
		return decodeUTF16(content[len(utf16LEBOM):], binary.LittleEndian)
	}
	if bytes.HasPrefix(content, utf16BEBOM) {
		return decodeUTF16(content[len(utf16BEBOM):], binary.BigEndian)
	}
	if utf8.Valid(content) {
		return string(content)
	}
	if decoded, err := simplifiedchinese.GB18030.NewDecoder().Bytes(content); err == nil && utf8.Valid(decoded) {
		return string(decoded)
	}
	return string(bytes.ToValidUTF8(content, []byte("\uFFFD")))
}

func decodeUTF16(content []byte, order binary.ByteOrder) string {
	if len(content)%2 == 1 {
		content = content[:len(content)-1]
	}
	u16 := make([]uint16, 0, len(content)/2)
	for len(content) >= 2 {
		u16 = append(u16, order.Uint16(content[:2]))
		content = content[2:]
	}
	return string(utf16.Decode(u16))
}
