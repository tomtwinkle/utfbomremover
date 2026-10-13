package utfbomremover_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"

	"golang.org/x/text/transform"

	"github.com/tomtwinkle/utfbomremover"
)

var edgeBOMs = []struct {
	name string
	data []byte
}{
	{name: "UTF-8", data: []byte{0xEF, 0xBB, 0xBF}},
	{name: "UTF-16BE", data: []byte{0xFE, 0xFF}},
	{name: "UTF-16LE", data: []byte{0xFF, 0xFE}},
	{name: "UTF-32BE", data: []byte{0x00, 0x00, 0xFE, 0xFF}},
	{name: "UTF-32LE", data: []byte{0xFF, 0xFE, 0x00, 0x00}},
}

func TestTransformerPartialSignatures(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
		want []byte
	}{
		{name: "UTF-32BE/one byte", data: []byte{0x00}, want: []byte{0x00}},
		{name: "UTF-32BE/two bytes", data: []byte{0x00, 0x00}, want: []byte{0x00, 0x00}},
		{name: "UTF-32BE/three bytes", data: []byte{0x00, 0x00, 0xFE}, want: []byte{0x00, 0x00, 0xFE}},
		{name: "UTF-32LE/one byte", data: []byte{0xFF}, want: []byte{0xFF}},
		// At EOF the complete UTF-16LE marker is removed, preserving any tail.
		{name: "UTF-32LE/two bytes", data: []byte{0xFF, 0xFE}},
		{name: "UTF-32LE/three bytes", data: []byte{0xFF, 0xFE, 0x00}, want: []byte{0x00}},
		{name: "UTF-8/one byte", data: []byte{0xEF}, want: []byte{0xEF}},
		{name: "UTF-8/two bytes", data: []byte{0xEF, 0xBB}, want: []byte{0xEF, 0xBB}},
		{name: "UTF-16BE/one byte", data: []byte{0xFE}, want: []byte{0xFE}},
	}
	marker := []byte{0xEF, 0xBB, 0xBF}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Run("EOF", func(t *testing.T) {
				assertTransformedStreams(t, testCase.data, testCase.want)
			})
			t.Run("EOF after marker", func(t *testing.T) {
				input := append(bytes.Clone(marker), testCase.data...)
				assertTransformedStreams(t, input, testCase.want)
			})
			t.Run("wait then mismatch", func(t *testing.T) {
				tr := utfbomremover.NewTransformer()
				input := append(bytes.Clone(marker), testCase.data...)
				assertTransformStep(t, tr, input, false, 16, nil, len(marker), transform.ErrShortSrc)
				assertTransformStep(t, tr, nil, false, 16, nil, 0, nil)
				// The caller resubmits the unconsumed prefix with the next byte.
				input = append(bytes.Clone(testCase.data), 'A')
				want := append(bytes.Clone(testCase.want), 'A')
				assertTransformStep(t, tr, input, false, 16, want, len(input), nil)
			})
		})
	}
}

func TestTransformerSplitSignatures(t *testing.T) {
	for _, bom := range edgeBOMs {
		for split := 1; split < len(bom.data); split++ {
			for _, prefix := range []struct {
				name string
				data []byte
			}{{name: "start"}, {name: "after marker", data: []byte{0xEF, 0xBB, 0xBF}}} {
				t.Run(fmt.Sprintf("%s/split-%d/%s", bom.name, split, prefix.name), func(t *testing.T) {
					tr := utfbomremover.NewTransformer()
					input := append(bytes.Clone(prefix.data), bom.data[:split]...)
					assertTransformStep(t, tr, input, false, 16, nil, len(prefix.data), transform.ErrShortSrc)
					input = append(bytes.Clone(bom.data), 'A')
					assertTransformStep(t, tr, input, false, 16, []byte{'A'}, len(input), nil)
					assertTransformStep(t, tr, nil, true, 16, nil, 0, nil)
				})
			}
		}
	}
}

func TestTransformerEmptyAndBOMOnly(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assertTransformedStreams(t, nil, nil)
	})
	t.Run("one content byte", func(t *testing.T) {
		assertTransformedStreams(t, []byte{'A'}, []byte{'A'})
	})
	for _, bom := range edgeBOMs {
		t.Run(bom.name, func(t *testing.T) {
			assertTransformedStreams(t, bytes.Repeat(bom.data, 2), nil)
		})
	}
}

func TestTransformerShortDestinations(t *testing.T) {
	markers := []byte{0x00, 0x00, 0xFE, 0xFF, 0xEF, 0xBB, 0xBF, 0xFE, 0xFF}
	payload := []byte{'A', 0xEF, 0xBB, 0xBF, 'B'}
	for _, size := range []int{1, 2} {
		t.Run(fmt.Sprintf("size-%d", size), func(t *testing.T) {
			tr := utfbomremover.NewTransformer()
			input := append(bytes.Clone(markers), payload...)
			// Even with no destination space, only the complete markers are consumed.
			assertTransformStep(t, tr, input, true, 0, nil, len(markers), transform.ErrShortDst)
			for remaining := payload; len(remaining) > 0; {
				want := remaining[:min(size, len(remaining))]
				var wantErr error
				if len(remaining) > size {
					wantErr = transform.ErrShortDst
				}
				assertTransformStep(t, tr, remaining, true, size, want, len(want), wantErr)
				remaining = remaining[len(want):]
			}
		})
	}
}

func TestTransformerResetStates(t *testing.T) {
	marker := []byte{0xEF, 0xBB, 0xBF}
	testCases := []struct {
		name  string
		src   []byte
		atEOF bool
		want  []byte
		nSrc  int
		err   error
	}{
		{name: "waiting for source", src: []byte{0xEF, 0xBB, 0xBF, 0xEF}, nSrc: 3, err: transform.ErrShortSrc},
		{name: "marker consumed", src: marker, nSrc: 3},
		{name: "copying content", src: []byte{'A'}, want: []byte{'A'}, nSrc: 1},
		{name: "empty EOF", atEOF: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tr := utfbomremover.NewTransformer()
			assertTransformStep(t, tr, testCase.src, testCase.atEOF, 16, testCase.want, testCase.nSrc, testCase.err)
			tr.Reset()
			tr.Reset()
			input := append(bytes.Repeat(marker, 2), 'A')
			assertTransformStep(t, tr, input, true, 16, []byte{'A'}, len(input), nil)
		})
	}
}

func assertTransformStep(t *testing.T, tr transform.Transformer, src []byte, atEOF bool, dstSize int, want []byte, wantSrc int, wantErr error) {
	t.Helper()
	dst := make([]byte, dstSize)
	nDst, nSrc, err := tr.Transform(dst, src, atEOF)
	if nDst != len(want) || nSrc != wantSrc || !errors.Is(err, wantErr) {
		t.Fatalf("Transform(% X, atEOF=%t, dst=%d) = (%d, %d, %v), want (%d, %d, %v)",
			src, atEOF, dstSize, nDst, nSrc, err, len(want), wantSrc, wantErr)
	}
	if !bytes.Equal(dst[:nDst], want) {
		t.Fatalf("transformed bytes = % X, want % X", dst[:nDst], want)
	}
}

func assertTransformedStreams(t *testing.T, input, want []byte) {
	t.Helper()
	assertTransformed(t, input, want)
	assertTransformedBytewise(t, input, want)
	actual, err := io.ReadAll(transform.NewReader(iotest.OneByteReader(bytes.NewReader(input)), utfbomremover.NewTransformer()))
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if !bytes.Equal(actual, want) {
		t.Fatalf("reader bytes = % X, want % X", actual, want)
	}
}
