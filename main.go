package utfbomremover

import (
	"bytes"

	"golang.org/x/text/transform"
)

func NewTransformer() transform.Transformer {
	return &remover{nop: transform.Nop}
}

type remover struct {
	nop  transform.Transformer
	done bool
}

var _ transform.Transformer = (*remover)(nil)

const (
	BOMSize4Byte = 4
	BOMSize3Byte = 3
	BOMSize2Byte = 2
)

var bomSignatures = [][]byte{
	{0x00, 0x00, 0xFE, 0xFF}, // UTF-32 Big Endian
	{0xFF, 0xFE, 0x00, 0x00}, // UTF-32 Little Endian
	{0xEF, 0xBB, 0xBF},       // UTF-8
	{0xFE, 0xFF},             // UTF-16 Big Endian
	{0xFF, 0xFE},             // UTF-16 Little Endian
}

func (t *remover) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if t.done {
		return t.nop.Transform(dst, src, atEOF)
	}

	if len(src) == 0 {
		if atEOF {
			t.done = true
		}
		return
	}

	for nSrc < len(src) {
		bomSize, partial := leadingBOMSize(src[nSrc:], atEOF)
		if partial {
			return 0, nSrc, transform.ErrShortSrc
		}
		if bomSize == 0 {
			t.done = true
			break
		}
		nSrc += bomSize
	}

	if nSrc == len(src) {
		if atEOF {
			t.done = true
		}
		return
	}

	remaining := src[nSrc:]
	nDst = copy(dst, remaining)
	nSrc += nDst
	if nDst < len(remaining) {
		err = transform.ErrShortDst
	}
	return
}

func (t *remover) Reset() {
	t.done = false
}

// leadingBOMSize reports the size of a BOM at the start of buf. If the bytes
// could still become a BOM, partial is true so the caller can wait for input.
func leadingBOMSize(buf []byte, atEOF bool) (size int, partial bool) {
	if len(buf) == 0 {
		return 0, false
	}

	var matchingSize int
	for _, signature := range bomSignatures {
		if bytes.HasPrefix(buf, signature) && len(signature) > matchingSize {
			matchingSize = len(signature)
		}

		if !atEOF && len(buf) < len(signature) && bytes.Equal(buf, signature[:len(buf)]) {
			return 0, true
		}
	}

	return matchingSize, false
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
