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

//nolint:gocyclo
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
			arrange: func(_ *testing.T) (Param, Want) {
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
			arrange: func(_ *testing.T) (Param, Want) {
				return Param{
					data: []byte{0x00, 0x00, 0xFE, 0xFF},
				}, Want{
					data: []byte{},
				}
			},
		},
		"UTF-32 LittleEndian:BOM only": {
			arrange: func(_ *testing.T) (Param, Want) {
				return Param{
					data: []byte{0xFF, 0xFE, 0x00, 0x00},
				}, Want{
					data: []byte{},
				}
			},
		},
		"UTF-8:BOM only": {
			arrange: func(_ *testing.T) (Param, Want) {
				return Param{
					data: []byte{0xEF, 0xBB, 0xBF},
				}, Want{
					data: []byte{},
				}
			},
		},
		"UTF-16 BigEndian:BOM only": {
			arrange: func(_ *testing.T) (Param, Want) {
				return Param{
					data: []byte{0xFE, 0xFF},
				}, Want{
					data: []byte{},
				}
			},
		},
		"UTF-16 LittleEndian:BOM only": {
			arrange: func(_ *testing.T) (Param, Want) {
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
			arrange: func(_ *testing.T) (Param, Want) {
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
		"UTF-8:multiple BOMs": {
			arrange: func(_ *testing.T) (Param, Want) {
				return Param{
					data: []byte{0xEF, 0xBB, 0xBF, 0xEF, 0xBB, 0xBF},
				}, Want{
					data: []byte{},
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
			if !bytes.Equal(want.data, actual.Bytes()) {
				t.Errorf("byte does not match\n%v", actual.Bytes())
			}
		})
	}
}

func TestTransformerMultipleLeadingBOMs(t *testing.T) {
	// Repeated markers are removed by policy, even though the same bytes can
	// represent an intentional leading U+FEFF character after the first marker.
	testCases := []struct {
		name    string
		bom     []byte
		payload []byte
	}{
		{name: "UTF-8", bom: []byte{0xEF, 0xBB, 0xBF}, payload: []byte("A")},
		{name: "UTF-16BE", bom: []byte{0xFE, 0xFF}, payload: []byte{0x00, 0x41}},
		{name: "UTF-16LE", bom: []byte{0xFF, 0xFE}, payload: []byte{0x41, 0x00}},
		{name: "UTF-32BE", bom: []byte{0x00, 0x00, 0xFE, 0xFF}, payload: []byte{0x00, 0x00, 0x00, 0x41}},
		{name: "UTF-32LE", bom: []byte{0xFF, 0xFE, 0x00, 0x00}, payload: []byte{0x41, 0x00, 0x00, 0x00}},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			for _, repetition := range []struct {
				count int
				name  string
			}{{count: 2, name: "twice"}, {count: 3, name: "three times"}} {
				repeatedBOMs := make([][]byte, repetition.count)
				for i := range repeatedBOMs {
					repeatedBOMs[i] = testCase.bom
				}
				input := append(bytes.Join(repeatedBOMs, nil), testCase.payload...)
				t.Run(repetition.name, func(t *testing.T) {
					assertTransformed(t, input, testCase.payload)
				})
			}

			input := append(bytes.Repeat(testCase.bom, 2), testCase.payload...)
			t.Run("bytewise writes", func(t *testing.T) {
				assertTransformedBytewise(t, input, testCase.payload)
			})
		})
	}

	t.Run("mixed encodings", func(t *testing.T) {
		markers := [][]byte{
			{0xEF, 0xBB, 0xBF},
			{0x00, 0x00, 0xFE, 0xFF},
			{0xFE, 0xFF},
			{0xFF, 0xFE},
			{0xFF, 0xFE, 0x00, 0x00},
			{0xEF, 0xBB, 0xBF},
		}
		input := append(bytes.Join(markers, nil), []byte("payload")...)
		assertTransformed(t, input, []byte("payload"))
	})

	t.Run("BOM after content is preserved", func(t *testing.T) {
		input := []byte{0xEF, 0xBB, 0xBF, 'A', 0xEF, 0xBB, 0xBF, 'B'}
		want := []byte{'A', 0xEF, 0xBB, 0xBF, 'B'}
		assertTransformed(t, input, want)
	})

	t.Run("near signatures are preserved", func(t *testing.T) {
		for _, input := range [][]byte{
			{0xEF, 0xBB, 0xBE},
			{0xFE, 0xFE},
			{0x00, 0x00, 0xFE, 0xFE},
		} {
			assertTransformed(t, input, input)
		}
	})

	t.Run("incomplete signature at EOF is preserved", func(t *testing.T) {
		for _, input := range [][]byte{
			{0xEF, 0xBB},
			{0x00, 0x00, 0xFE},
		} {
			assertTransformed(t, input, input)
		}
	})

	t.Run("UTF-32LE signature overlaps UTF-16LE BOM and NUL", func(t *testing.T) {
		// FF FE 00 00 is either a UTF-32LE BOM or a UTF-16LE BOM followed by U+0000.
		input := []byte{0xFF, 0xFE, 0x00, 0x00}
		assertTransformed(t, input, nil)
	})
}

func assertTransformed(t *testing.T, input, want []byte) {
	t.Helper()
	actual, nSrc, err := transform.Bytes(utfbomremover.NewTransformer(), input)
	if err != nil {
		t.Fatalf("transform.Bytes() error = %v", err)
	}
	if nSrc != len(input) {
		t.Fatalf("consumed source bytes = %d, want %d", nSrc, len(input))
	}
	if !bytes.Equal(actual, want) {
		t.Fatalf("transformed bytes = % X, want % X", actual, want)
	}
}

func assertTransformedBytewise(t *testing.T, input, want []byte) {
	t.Helper()
	var output bytes.Buffer
	writer := transform.NewWriter(&output, utfbomremover.NewTransformer())
	for i := range input {
		if n, err := writer.Write(input[i : i+1]); err != nil {
			t.Fatalf("Write() error = %v", err)
		} else if n != 1 {
			t.Fatalf("Write() consumed %d bytes, want 1", n)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !bytes.Equal(output.Bytes(), want) {
		t.Fatalf("transformed bytes = % X, want % X", output.Bytes(), want)
	}
}

// nolint: typecheck, gocyclo
func FuzzTransformer(f *testing.F) {
	bomseeds := [][]byte{
		{},
		{0x00, 0x00, 0xFE, 0xFF},
		{0x00, 0x00, 0xFE, 0xFF, 0x00, 0x00, 0xFE, 0xFF},
		{0xFF, 0xFE, 0x00, 0x00},
		{0xFF, 0xFE, 0x00, 0x00, 0xFF, 0xFE, 0x00, 0x00},
		{0xEF, 0xBB, 0xBF},
		{0xEF, 0xBB, 0xBF, 0xEF, 0xBB, 0xBF},
		{0xFE, 0xFF},
		{0xFE, 0xFF, 0xFE, 0xFF},
		{0xFF, 0xFE},
		{0xFF, 0xFE, 0xFF, 0xFE},
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
		tr := utfbomremover.NewTransformer()
		for len(p) > 0 {
			r, n, err := transform.Bytes(tr, p)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if utfbomremover.IsUTF32BigEndianBOM(r) {
				t.Fatalf("utf-32 be bom [%X]", r)
			}
			if utfbomremover.IsUTF32LittleEndianBOM(r) {
				t.Fatalf("utf-32 le bom [%X]", r)
			}
			if utfbomremover.IsUTF8BOM(r) {
				t.Fatalf("utf-8 bom [%X]", r)
			}
			if utfbomremover.IsUTF16BigEndianBOM(r) {
				t.Fatalf("utf-16 be bom [%X]", r)
			}
			if utfbomremover.IsUTF16LittleEndianBOM(r) {
				t.Fatalf("utf-16 le bom [%X]", r)
			}
			p = p[n:]
			tr.Reset()
		}
	})
}
