package miio

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"errors"
	"fmt"
)

const HelloPacketDeviceId uint32 = 0xFFFFFFFF

var (
	ErrPadding             = errors.New("Padding size error while decrypting payload")
	ErrCryptoNotSet        = errors.New("Crypto key/iv not set")
	ErrWrongPacket         = errors.New("Wrong packet, magic is illegal")
	ErrInvalidPacketType   = errors.New("invalid packet type")
	ErrInvalidPacketLength = errors.New("invalid packet Len")
	ErrUnknownPacket       = errors.New("unknown packet type")
	ErrReadFromBuf         = errors.New("error read data from buffer")
)

type Packet struct {
	// crypto key/iv
	iv          []byte
	key         []byte
	DeviceToken []byte

	// packet header
	Magic     uint16
	Length    uint16
	Unknown1  uint32
	DeviceId  uint32
	Timestamp uint32
	CheckSum  []byte // 16

	// payload
	Data []byte
}

func NewPacket(deviceId uint32, deviceToken []byte, timestamp uint32, payload []byte) (*Packet, error) {
	packet := &Packet{
		Magic:     0x2131,
		Length:    uint16(0x20),
		Unknown1:  0x0,
		DeviceId:  deviceId,
		Timestamp: timestamp,
	}

	if deviceId == HelloPacketDeviceId {
		packet.Unknown1 = HelloPacketDeviceId
		packet.CheckSum = bytes.Repeat([]byte{0xFF}, 16)
	} else if deviceToken != nil {
		packet.CheckSum = deviceToken
		packet.DeviceToken = deviceToken

		hash := md5.New()
		if _, err := hash.Write(deviceToken); err != nil {
			return nil, err
		}
		packet.key = hash.Sum(nil)

		hash = md5.New()
		if _, err := hash.Write(packet.key); err != nil {
			return nil, err
		}
		if _, err := hash.Write(deviceToken); err != nil {
			return nil, err
		}
		packet.iv = hash.Sum(nil)

		if payload != nil {
			if err := packet.Encrypt(payload); err != nil {
				return nil, err
			}
		}
	}

	return packet, nil
}

func ParsePacket(deviceId uint32, deviceToken []byte, buf []byte) (*Packet, error) {
	packet := &Packet{
		DeviceId:    deviceId,
		DeviceToken: deviceToken,
	}

	err := packet.Unpack(buf)
	if err != nil {
		return nil, err
	}

	if deviceToken != nil {
		hash := md5.New()
		_, err = hash.Write(deviceToken)
		if err != nil {
			return nil, err
		}
		packet.key = hash.Sum(nil)

		hash = md5.New()
		_, err = hash.Write(packet.key)
		if err != nil {
			return nil, err
		}
		_, err = hash.Write(deviceToken)
		if err != nil {
			return nil, err
		}
		packet.iv = hash.Sum(nil)

		packet.Data, err = packet.Decrypt()
		if err != nil {
			return nil, err
		}
	}

	return packet, nil
}

func (p *Packet) pkcs5Pad(data []byte, blockSize int) []byte {
	length := len(data)
	padLength := (blockSize - (length % blockSize))
	pad := bytes.Repeat([]byte{byte(padLength)}, padLength)
	return append(data, pad...)
}

func (p *Packet) pkcs5Unpad(data []byte, blockSize int) ([]byte, error) {
	srcLen := len(data)
	paddingLen := int(data[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, ErrPadding
	}
	return data[:srcLen-paddingLen], nil
}

func (p *Packet) Pack() ([]byte, error) {
	var buf []byte = make([]byte, 16)
	var offset int = 0
	offset = WriteInt16(buf, offset, p.Magic)
	offset = WriteInt16(buf, offset, p.Length)
	offset = WriteInt32(buf, offset, p.Unknown1)
	offset = WriteInt32(buf, offset, p.DeviceId)
	offset = WriteInt32(buf, offset, p.Timestamp)

	if p.DeviceId == HelloPacketDeviceId {
		buf = append(buf, p.CheckSum...)
		return buf, nil
	}

	b := append(buf, p.CheckSum...)
	b = append(b, p.Data...)

	// calculate checksum
	h := md5.New()
	h.Write(b)

	// ...  put into structure
	p.CheckSum = h.Sum(nil)
	buf = append(buf, p.CheckSum...)
	buf = append(buf, p.Data...)

	return buf, nil
}

func (p *Packet) Unpack(buf []byte) error {
	var offset int = 0
	var err error

	if p.Magic, offset, err = ReadInt16(buf, offset); err != nil {
		return err
	}
	if p.Magic != 0x2131 {
		return ErrWrongPacket
	}
	if p.Length, offset, err = ReadInt16(buf, offset); err != nil {
		return err
	}
	if p.Unknown1, offset, err = ReadInt32(buf, offset); err != nil {
		return err
	}
	if p.DeviceId, offset, err = ReadInt32(buf, offset); err != nil {
		return err
	}
	if p.Timestamp, offset, err = ReadInt32(buf, offset); err != nil {
		return err
	}
	if p.CheckSum, offset, err = ReadBytes(buf, offset, 16); err != nil {
		return err
	}

	if p.Length-0x20 > 0 {
		if p.Data, offset, err = ReadBytes(buf, offset, int(p.Length)-0x20); err != nil {
			return err
		}
	}

	return nil
}

func (p *Packet) Decrypt() ([]byte, error) {
	if p.key == nil || p.iv == nil {
		return nil, ErrCryptoNotSet
	}

	block, err := aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCBCDecrypter(block, p.iv)
	decrypted := make([]byte, len(p.Data))
	stream.CryptBlocks(decrypted, p.Data)

	return p.pkcs5Unpad(decrypted, block.BlockSize())
}

func (p *Packet) Encrypt(payload []byte) error {
	if p.key == nil || p.iv == nil {
		return ErrCryptoNotSet
	}

	block, err := aes.NewCipher(p.key)
	if err != nil {
		return err
	}

	data := p.pkcs5Pad(payload, block.BlockSize())
	stream := cipher.NewCBCEncrypter(block, p.iv)

	encrypted := make([]byte, len(data))
	stream.CryptBlocks(encrypted, data)

	p.Data = encrypted
	p.Length += uint16(len(encrypted))

	return nil
}

func (p *Packet) String() string {
	return fmt.Sprintf(`{"magic":"%x","len":%d,"unknown1":"%x","id":"%x","token":"%x","timestamp":"%d","checksum":"%x","data":[%s]}`,
		p.Magic, p.Length, p.Unknown1, p.DeviceId, p.DeviceToken, p.Timestamp, p.CheckSum, p.Data)
}
