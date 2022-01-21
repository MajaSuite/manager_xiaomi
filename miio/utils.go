package miio

import (
	"bytes"
	"encoding/binary"
)

func WriteInt8(buf []byte, offset int, value uint8) int {
	buf[offset] = value
	offset++
	return offset
}

func WriteInt16(buf []byte, offset int, value uint16) int {
	binary.BigEndian.PutUint16(buf[offset:], value)
	offset += 2
	return offset
}

func WriteInt32(buf []byte, offset int, value uint32) int {
	binary.BigEndian.PutUint32(buf[offset:], value)
	offset += 4
	return offset
}

func ReadInt8(buf []byte, offset int) (uint8, int, error) {
	var value uint8

	if len(buf) >= (offset + 1) {
		value = buf[offset]
		return value, offset + 1, nil
	}

	return 0, offset, ErrReadFromBuf
}

func ReadInt16(buf []byte, offset int) (uint16, int, error) {
	var value uint16

	if len(buf) >= (offset + 2) {
		value = uint16(buf[offset+1]) | uint16(buf[offset])<<8
		return value, offset + 2, nil
	}

	return 0, offset, ErrReadFromBuf
}

func ReadInt32(buf []byte, offset int) (uint32, int, error) {
	var value uint32

	if len(buf) >= (offset + 4) {
		value = uint32(buf[offset+3]) |
			uint32(buf[offset+2])<<8 |
			uint32(buf[offset+1])<<16 |
			uint32(buf[offset])<<24

		return value, offset + 4, nil
	}

	return 0, offset, ErrReadFromBuf
}

func ReadBytes(buf []byte, offset int, length int) ([]byte, int, error) {
	var value []byte

	if len(buf) >= (offset + length) {
		value = buf[offset : offset+length]
		return value, offset + length, nil
	}

	return nil, offset, ErrReadFromBuf
}

func ReadString(buf []byte, offset int, length int) (string, int, error) {
	s, i, e := ReadBytes(buf, offset, length)
	return string(s), i, e
}

// Pad using PKCS5 padding scheme.
func pkcs5Pad(data []byte, blockSize int) []byte {
	length := len(data)
	padLength := (blockSize - (length % blockSize))
	pad := bytes.Repeat([]byte{byte(padLength)}, padLength)
	return append(data, pad...)
}

// Unpad using PKCS5 padding scheme.
func pkcs5Unpad(data []byte, blockSize int) ([]byte, error) {
	srcLen := len(data)
	paddingLen := int(data[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, ErrPadding
	}
	return data[:srcLen-paddingLen], nil
}
