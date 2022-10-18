package utfbomremover

import (
	"golang.org/x/text/transform"
)

func NewTransformer() transform.Transformer {
	return &remover{nop: transform.Nop}
}

type remover struct {
	nop     transform.Transformer
	counter int
}

var _ transform.Transformer = (*remover)(nil)

const (
	BOMSize4Byte = 4
	BOMSize3Byte = 3
	BOMSize2Byte = 2
)

func (t *remover) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	_src := src
	if len(_src) == 0 && atEOF {
		return
	}

	if t.counter > 0 {
		// BOMは先頭にしか存在しないため1回目以外はSpanningTransformerに委譲する
		return t.nop.Transform(dst, src, atEOF)
	}

	// TODO: buffer size が BOM 以下の場合おそらくBOM削除されないので要確認
	var (
		buf         []byte
		writeBufLen int
		remainder   int
	)
	if len(_src) >= len(dst) {
		buf = _src[:len(dst)]
		remainder = len(_src) - len(buf)
	} else {
		buf = _src
	}

	writeBufLen = len(buf)
	switch {
	case IsUTF32BigEndianBOM(buf):
		buf = buf[BOMSize4Byte:]
	case IsUTF32LittleEndianBOM(buf):
		buf = buf[BOMSize4Byte:]
	case IsUTF8BOM(buf):
		buf = buf[BOMSize3Byte:]
	case IsUTF16BigEndianBOM(buf):
		buf = buf[BOMSize2Byte:]
	case IsUTF16LittleEndianBOM(buf):
		buf = buf[BOMSize2Byte:]
	}

	dstN := copy(dst, buf)
	nSrc = writeBufLen
	nDst = dstN
	if remainder > 0 {
		// over destination buffer
		err = transform.ErrShortDst
	}
	t.counter++
	return
}

func (t *remover) Reset() {
	t.counter = 0
}

func IsUTF32BigEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize4Byte && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0xFE && buf[3] == 0xFF
}

func IsUTF32LittleEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize4Byte && buf[0] == 0xFF && buf[1] == 0xFE && buf[2] == 0x00 && buf[3] == 0x00
}

func IsUTF8BOM(buf []byte) bool {
	return len(buf) >= BOMSize3Byte && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF
}

func IsUTF16BigEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize2Byte && buf[0] == 0xFE && buf[1] == 0xFF
}

func IsUTF16LittleEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize2Byte && buf[0] == 0xFF && buf[1] == 0xFE
}
