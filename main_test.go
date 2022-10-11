package utfbomremover_test

import (
	"bytes"
	"errors"
	"testing"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"

	"github.com/tomtwinkle/utfbomremover"
)

func TestNewTransformer(t *testing.T) {
	type Param struct {
		data []byte
	}
	type Want struct {
		data []byte
	}

	tests := map[string]struct {
		arrange   func(*testing.T) (Param, Want)
		wantError error
	}{
		"UTF-32 BigEndian:no BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-32 LittleEndian:no BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-8:no BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				data := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				return Param{
						data: data,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-16 BigEndian:no BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-16 LittleEndian:no BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-32 BigEndian:BOM only": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0x00, 0x00, 0xFE, 0xFF},
					}, Want{
						data: []byte{},
					}
			},
		},
		"UTF-32 LittleEndian:BOM only": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0xFF, 0xFE, 0x00, 0x00},
					}, Want{
						data: []byte{},
					}
			},
		},
		"UTF-8:BOM only": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0xEF, 0xBB, 0xBF},
					}, Want{
						data: []byte{},
					}
			},
		},
		"UTF-16 BigEndian:BOM only": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0xFE, 0xFF},
					}, Want{
						data: []byte{},
					}
			},
		},
		"UTF-16 LittleEndian:BOM only": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0xFF, 0xFE},
					}, Want{
						data: []byte{},
					}
			},
		},
		"UTF-32 BigEndian:BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := utf32.UTF32(utf32.BigEndian, utf32.UseBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data[utfbomremover.BOMSize4Byte:],
					}
			},
		},
		"UTF-32 LittleEndian:BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := utf32.UTF32(utf32.LittleEndian, utf32.UseBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data[utfbomremover.BOMSize4Byte:],
					}
			},
		},
		"UTF-8:BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				data := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				bomData := append([]byte{0xEF, 0xBB, 0xBF}, data...)
				return Param{
						data: bomData,
					}, Want{
						data: data,
					}
			},
		},
		"UTF-16 BigEndian:BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data[utfbomremover.BOMSize2Byte:],
					}
			},
		},
		"UTF-16 LittleEndian:BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				e := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()
				org := bytes.Repeat([]byte("abcd一二三四五六七八九十拾壱"), 1000)
				data, err := e.Bytes(org)
				if err != nil {
					t.Fatal(err)
				}
				return Param{
						data: data,
					}, Want{
						data: data[utfbomremover.BOMSize2Byte:],
					}
			},
		},
		"Illegal case UTF-8:BOM BOM": {
			arrange: func(t *testing.T) (Param, Want) {
				return Param{
						data: []byte{0xEF, 0xBB, 0xBF, 0xEF, 0xBB, 0xBF},
					}, Want{
						data: []byte{0xEF, 0xBB, 0xBF},
					}
			},
		},
	}

	for n, v := range tests {
		name := n
		tt := v

		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			param, want := tt.arrange(t)
			w := transform.NewWriter(&buf, utfbomremover.NewTransformer())
			if _, err := w.Write(param.data); err != nil {
				if tt.wantError != nil && errors.Is(err, tt.wantError) {
					return
				}
				t.Error(err)
			}
			if err := w.Close(); err != nil {
				t.Error(err)
			}
			var actual bytes.Buffer
			if _, err := actual.Write(buf.Bytes()); err != nil {
				t.Error(err)
			}

			if len(want.data) != actual.Len() {
				t.Errorf("byte length does not match %d=%d", len(want.data), actual.Len())
			}
			if string(want.data) != string(actual.Bytes()) {
				t.Errorf("byte does not match\n%v", actual.Bytes())
			}
		})
	}
}

// nolint: typecheck
func FuzzTransformer(f *testing.F) {
	bomseeds := [][]byte{
		{0x00, 0x00, 0xFE, 0xFF},
		{0xFF, 0xFE, 0x00, 0x00},
		{0xEF, 0xBB, 0xBF},
		{0xFE, 0xFF},
		{0xFF, 0xFE},
	}
	seeds := [][]byte{
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("一二三四五六七八九十拾壱"),
		[]byte("咖呸咕咀呻呷咄咒咆呼咐呱呶和咚呢"),
		[]byte("🍣🍺🍥🍜💯"),
	}

	for _, b := range append(bomseeds, seeds...) {
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, p []byte) {
		if (utfbomremover.ISUTF32BigEndianBOM(p) && utfbomremover.ISUTF32BigEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32BigEndianBOM(p) && utfbomremover.ISUTF32LittleEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32BigEndianBOM(p) && utfbomremover.ISUTF8BOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32BigEndianBOM(p) && utfbomremover.ISUTF16BigEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32BigEndianBOM(p) && utfbomremover.ISUTF16LittleEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32LittleEndianBOM(p) && utfbomremover.ISUTF32BigEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32LittleEndianBOM(p) && utfbomremover.ISUTF32LittleEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32LittleEndianBOM(p) && utfbomremover.ISUTF8BOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32LittleEndianBOM(p) && utfbomremover.ISUTF16BigEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF32LittleEndianBOM(p) && utfbomremover.ISUTF16LittleEndianBOM(p[utfbomremover.BOMSize4Byte:])) ||
			(utfbomremover.ISUTF8BOM(p) && utfbomremover.ISUTF32BigEndianBOM(p[utfbomremover.BOMSize3Byte:])) ||
			(utfbomremover.ISUTF8BOM(p) && utfbomremover.ISUTF32LittleEndianBOM(p[utfbomremover.BOMSize3Byte:])) ||
			(utfbomremover.ISUTF8BOM(p) && utfbomremover.ISUTF8BOM(p[utfbomremover.BOMSize3Byte:])) ||
			(utfbomremover.ISUTF8BOM(p) && utfbomremover.ISUTF16BigEndianBOM(p[utfbomremover.BOMSize3Byte:])) ||
			(utfbomremover.ISUTF8BOM(p) && utfbomremover.ISUTF16LittleEndianBOM(p[utfbomremover.BOMSize3Byte:])) ||
			(utfbomremover.ISUTF16BigEndianBOM(p) && utfbomremover.ISUTF32BigEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16BigEndianBOM(p) && utfbomremover.ISUTF32LittleEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16BigEndianBOM(p) && utfbomremover.ISUTF8BOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16BigEndianBOM(p) && utfbomremover.ISUTF16BigEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16BigEndianBOM(p) && utfbomremover.ISUTF16LittleEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16LittleEndianBOM(p) && utfbomremover.ISUTF32BigEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16LittleEndianBOM(p) && utfbomremover.ISUTF32LittleEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16LittleEndianBOM(p) && utfbomremover.ISUTF8BOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16LittleEndianBOM(p) && utfbomremover.ISUTF16BigEndianBOM(p[utfbomremover.BOMSize2Byte:])) ||
			(utfbomremover.ISUTF16LittleEndianBOM(p) && utfbomremover.ISUTF16LittleEndianBOM(p[utfbomremover.BOMSize2Byte:])) {
			t.Skip()
		}
		tr := utfbomremover.NewTransformer()
		for len(p) > 0 {
			r, n, err := transform.Bytes(tr, p)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if utfbomremover.ISUTF32BigEndianBOM(r) {
				t.Fatalf("utf-32 be bom [%X]", r)
			}
			if utfbomremover.ISUTF32LittleEndianBOM(r) {
				t.Fatalf("utf-32 le bom [%X]", r)
			}
			if utfbomremover.ISUTF8BOM(r) {
				t.Fatalf("utf-8 bom [%X]", r)
			}
			if utfbomremover.ISUTF16BigEndianBOM(r) {
				t.Fatalf("utf-16 be bom [%X]", r)
			}
			if utfbomremover.ISUTF16LittleEndianBOM(r) {
				t.Fatalf("utf-16 le bom [%X]", r)
			}
			p = p[n:]
			tr.Reset()
		}
	})
}
