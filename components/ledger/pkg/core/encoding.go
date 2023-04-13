package core

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

type marshaller interface {
	Marshal(buf *Buffer)
	Unmarshal(buf *Buffer)
}

type Buffer struct {
	buf *bytes.Buffer
}

func (b *Buffer) must(err error) {
	if err != nil {
		panic(err)
	}
}

func (b *Buffer) writeByte(c byte) {
	b.must(b.buf.WriteByte(c))
}

func (b *Buffer) Reset() {
	b.buf.Reset()
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) write(bytes []byte) {
	must(b.buf.Write(bytes))
}

func (b *Buffer) writeUInt64(v uint64) {
	b.writeByte(byte(v >> 56))
	b.writeByte(byte(v >> 48))
	b.writeByte(byte(v >> 40))
	b.writeByte(byte(v >> 32))
	b.writeByte(byte(v >> 24))
	b.writeByte(byte(v >> 16))
	b.writeByte(byte(v >> 8))
	b.writeByte(byte(v))
}

func (b *Buffer) readUInt64() uint64 {
	return binary.BigEndian.Uint64(b.buf.Next(8))
}

func (b *Buffer) writeString(v string) {
	b.writeUInt64(uint64(len(v)))
	must(b.buf.WriteString(v))
}

func (b *Buffer) readString() string {
	size := b.readUInt64()
	return string(b.buf.Next(int(size)))
}

func (b *Buffer) writeDate(timestamp Time) {
	bytes := must(timestamp.MarshalBinary())
	b.writeUInt64(uint64(len(bytes)))
	b.write(bytes)
}

func (b *Buffer) readDate() Time {
	size := b.readUInt64()
	t := Time{}
	bytes := b.buf.Next(int(size))
	b.must(t.UnmarshalBinary(bytes))
	return t
}

func (b *Buffer) writeBigint(amount *big.Int) {
	bytes := amount.Bytes()
	b.writeUInt64(uint64(len(bytes)))
	b.write(bytes)
}

func (b *Buffer) readBigint() *big.Int {
	size := b.readUInt64()
	bytes := b.buf.Next(int(size))
	return (&big.Int{}).SetBytes(bytes)
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		buf: bytes.NewBuffer(data),
	}
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}

	return v
}
