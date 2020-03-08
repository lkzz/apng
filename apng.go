package apng

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// APNG 是基于 PNG 格式扩展的, 在PNG基础上引入了三个控制块
// 分别是：acTL（动画控制块）、fcTL（帧控制块）、fdAT（帧数据块）
// APNG 图片header格式如下：
// ------------------------------------------
// |  PNG Signature | IHDR chunk | acTL chunk|
// -------------------------------------------
// PNG Signature：是png图片固有的签名块,8个字节: 89 50 4e 47 0d 0a 1a 0a
// PNG chunk 格式：4字节chunk data length + 4字节chunk data type + chunk data + 4字节 crc
// IHDR: 图像头部块,格式：
// acTL: 动画控制块, 包含 字节帧数、显示次数以及4个固定字节固定：acTL

// chunk markers
var (
	pngHeader     = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	acTL          = []byte{0x61, 0x63, 0x54, 0x4c}
	IHDR          = []byte{0x49, 0x48, 0x44, 0x52}
	chkLen    int = 4
	typLen    int = 4
	crcLen    int = 4
)

type Decoder struct {
	buf *bytes.Buffer
}

type chunk struct {
	length uint32
	typ    []byte
	data   []byte
	crc    uint32
}

func NewDecoder(in []byte) (*Decoder, error) {
	d := &Decoder{buf: bytes.NewBuffer(in)}
	header := d.buf.Next(len(pngHeader))
	if !bytes.Equal(header, pngHeader) {
		return nil, fmt.Errorf("input is not a png file")
	}
	return d, nil
}

func (d *Decoder) ReadOneChunk() (*chunk, error) {
	ck := new(chunk)
	if err := binary.Read(d.buf, binary.BigEndian, &ck.length); err != nil {
		return nil, err
	}
	data := make([]byte, typLen+int(ck.length))
	if _, err := d.buf.Read(data); err != nil {
		return nil, err
	}
	if err := binary.Read(d.buf, binary.BigEndian, &ck.crc); err != nil {
		return nil, err
	}
	if crc := crc32.ChecksumIEEE(data); crc != ck.crc {
		return nil, fmt.Errorf("failed to check crc, expect:%d, get:%d", ck.crc, crc)
	}
	ck.typ = data[:typLen]
	ck.data = data[typLen:]
	return ck, nil
}

// Hit check input is apng or not.
func Hit(in []byte) bool {
	d, err := NewDecoder(in)
	if err != nil {
		return false
	}
	idhrChunk, err := d.ReadOneChunk()
	if err != nil || !bytes.Equal(idhrChunk.typ, IHDR) {
		return false
	}
	actlChunk, err := d.ReadOneChunk()
	if err != nil || !bytes.Equal(actlChunk.typ, acTL) {
		return false
	}
	return true
}
