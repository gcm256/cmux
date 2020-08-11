package cmux

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

const sizeInfoInBytes = 4
const startChar = ':'
const endChar = ','

// ParseNetString Finds out if a byte stream is a valid NetString
func ParseNetString(r *bufio.Reader) ([]byte, error) {
	strLen, err := parsNetStringLength(r)
	if err != nil {
		return []byte{}, err
	}
	if err = parseNetStringStartChar(r); err != nil {
		return []byte{}, err
	}
	b, err := parseNetStringInBytes(r, strLen)
	if err != nil {
		return []byte{}, err
	}
	if err = parserNetStringEndChar(r); err != nil {
		return []byte{}, err
	}
	return b, nil
}

// EncodeBytesToNetString Encodes a plain byte array to a netstring byte array
func EncodeBytesToNetString(data []byte) []byte {
	length := strconv.FormatInt(int64(len(data)), 10)
	return []byte(length + ":" + string(data) + ",")
}

func parsNetStringLength(r *bufio.Reader) (int, error) {
	buffer := make([]byte, sizeInfoInBytes)
	if _, err := io.ReadFull(r, buffer); err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(buffer)), nil
}

func parseNetStringInBytes(r *bufio.Reader, len int) ([]byte, error) {
	return ioutil.ReadAll(io.LimitReader(r, int64(len)))
}

func parseNetStringStartChar(r *bufio.Reader) error {
	return parseNetStringDelimChar(r, startChar)
}

func parserNetStringEndChar(r *bufio.Reader) error {
	return parseNetStringDelimChar(r, endChar)
}

func parseNetStringDelimChar(r *bufio.Reader, delimChar byte) error {
	char, err := r.ReadByte()
	if err != nil || char != delimChar {
		return fmt.Errorf("ERROR: Unexpected NetString delimiter char %c. Exppected %c", char, delimChar)
		//return err
	}
	return nil
}
