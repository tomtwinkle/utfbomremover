package utfbomremover

import (
	"golang.org/x/text/transform"
)

func NewTransformer() transform.Transformer {
	return &remover{}
}

type remover struct {
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
	if t.counter == 0 {
		switch {
		case ISUTF32BigEndianBOM(buf):
			buf = buf[BOMSize4Byte:]
		case ISUTF32LittleEndianBOM(buf):
			buf = buf[BOMSize4Byte:]
		case ISUTF8BOM(buf):
			buf = buf[BOMSize3Byte:]
		case ISUTF16BigEndianBOM(buf):
			buf = buf[BOMSize2Byte:]
		case ISUTF16LittleEndianBOM(buf):
			buf = buf[BOMSize2Byte:]
		}
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

func ISUTF32BigEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize4Byte && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0xFE && buf[3] == 0xFF
}

func ISUTF32LittleEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize4Byte && buf[0] == 0xFF && buf[1] == 0xFE && buf[2] == 0x00 && buf[3] == 0x00
}

func ISUTF8BOM(buf []byte) bool {
	return len(buf) >= BOMSize3Byte && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF
}

func ISUTF16BigEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize2Byte && buf[0] == 0xFE && buf[1] == 0xFF
}

func ISUTF16LittleEndianBOM(buf []byte) bool {
	return len(buf) >= BOMSize2Byte && buf[0] == 0xFF && buf[1] == 0xFE
}
