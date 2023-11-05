package packet

import "encoding/binary"

/*
   @Author: orbit-w
   @File: packet
   @2023 11月 周日 14:48
*/

type IPacket interface {
	Len() int
	Cap() int
	Off() uint
	Remain() []byte
	Data() []byte

	//writer
	Write(v []byte)
	WriteBool(v bool)
	WriteBytes(v []byte)
	WriteString(v string)
	WriteUint8(v uint8)
	WriteUint16(v uint16)
	WriteUint32(v uint32)
	WriteUint64(v uint64)

	//reader
	ReadBool() (ret bool, err error)
	ReadBytes() (ret []byte, err error)
	ReadBytes32() (ret []byte, err error)
	ReadInt8() (ret int8, err error)
	ReadUint16() (ret uint16, err error)
	ReadUint32() (ret uint32, err error)
	ReadUint64() (ret uint64, err error)
	NextBytesSize() (int, error)
	NextBytesSize32() (int, error)

	Reset()
	Return()
}

type Packet struct {
	off uint // read at &buf[off], write at &buf[len(buf)]
	buf []byte
}

func New() *Packet {
	return &Packet{
		buf: make([]byte, 0),
	}
}

func NewWithInitialSize(initSize int) *Packet {
	return &Packet{
		buf: make([]byte, initSize),
	}
}

func getPacket() *Packet {
	pack := pool.Get().(*Packet)
	return pack
}

func (p *Packet) Remain() []byte {
	return p.buf[p.off:]
}

func (p *Packet) Reset() {
	p.off = 0
	p.buf = p.buf[:0]
}

func (p *Packet) Data() []byte {
	return p.buf
}

func (p *Packet) Len() int {
	return len(p.buf)
}

func (p *Packet) Cap() int {
	return cap(p.buf)
}

func (p *Packet) Off() uint {
	return p.off
}

func (p *Packet) OutOfRange(n uint) bool {
	return p.off+n > uint(p.Len())
}

func (p *Packet) NextBytesSize() (int, error) {
	if p.OutOfRange(2) {
		return 0, ErrReadBytesHeaderFailed
	}
	buf := p.buf[p.off : p.off+2]
	return int(uint16(buf[0])<<8 | uint16(buf[1])), nil
}

func (p *Packet) NextBytesSize32() (int, error) {
	if p.OutOfRange(4) {
		return 0, ErrReadBytesHeaderFailed
	}
	buf := p.buf[p.off : p.off+4]
	return int(binary.BigEndian.Uint32(buf)), nil
}

func (p *Packet) Return() {
	p.Reset()
	pool.Put(p)
}
