package packet

import (
	"encoding/binary"
)

/*
   @Author: orbit-w
   @File: reader
   @2023 11月 周日 15:05
*/

func Reader(data []byte) IPacket {
	packet := getPacket()
	packet.buf = append(packet.buf, data...)
	return packet
}

func (p *Packet) ReadBool() (ret bool, err error) {
	var b byte
	b, err = p.ReadByte()
	if err != nil {
		return
	}

	return b == byte(1), nil
}

func (p *Packet) ReadByte() (ret byte, err error) {
	if p.off >= uint(p.Len()) {
		err = ErrReadByteFailed
		return
	}
	ret = p.buf[p.off]
	p.off++
	return
}

func (p *Packet) ReadInt8() (ret int8, err error) {
	ret = int8(p.buf[p.off])
	p.off++
	return
}

func (p *Packet) ReadUint16() (ret uint16, err error) {
	var shift uint = 2
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}

	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint16(buf)
	p.off += shift
	return
}

func (p *Packet) ReadUint32() (ret uint32, err error) {
	var shift uint = 4
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}

	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint32(buf)
	p.off += shift
	return
}

func (p *Packet) ReadUint64() (ret uint64, err error) {
	var shift uint = 8
	if p.OutOfRange(shift) {
		return 0, ErrOutOfRange
	}
	buf := p.buf[p.off : p.off+shift]
	ret = binary.BigEndian.Uint64(buf)
	p.off += shift
	return
}

func (p *Packet) ReadBytes() (ret []byte, err error) {
	v, rErr := p.ReadUint16()
	if rErr != nil {
		err = rErr
		return
	}

	shift := uint(v)
	ret = p.buf[p.off : p.off+shift]
	p.off += shift
	return
}

func (p *Packet) ReadBytes32() (ret []byte, err error) {
	v, rErr := p.ReadUint32()
	if rErr != nil {
		err = rErr
		return
	}

	shift := uint(v)
	ret = p.buf[p.off : p.off+shift]
	p.off += shift
	return
}
