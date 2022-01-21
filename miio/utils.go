package packet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

var ErrInvalidPacketType = errors.New("invalid packet type")
var ErrInvalidPacketLength = errors.New("invalid packet Len")
var ErrUnknownPacket = errors.New("unknown packet type")
var ErrReadFromBuf = errors.New("error read data from buffer")

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

func ReadPacket(conn net.Conn) ([]byte, error) {
	var magic, size uint16
	var err error

	reader := bufio.NewReader(conn)

	header := make([]byte, 4)
	if _, err := reader.Read(header); err != nil {
		return nil, err
	}

	if magic, _, err = ReadInt16(header, 0); err != nil {
		return nil, err
	}
	if magic != 0x2131 {
		return nil, errors.New("wrong magic " + strconv.Itoa(int(magic)))
	}

	if size, _, err = ReadInt16(header, 2); err != nil {
		return nil, err
	}
	res := make([]byte, size)
	WriteInt16(res, 0, magic)
	WriteInt16(res, 2, size)
	if _, err = reader.Read(res[4:]); err != nil {
		return nil, err
	}

	return res, nil
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
