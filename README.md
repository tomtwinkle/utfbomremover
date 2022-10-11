# utfbomremover

![GitHub](https://img.shields.io/github/license/tomtwinkle/utfbomremover)
[![Go Report Card](https://goreportcard.com/badge/github.com/olvrng/ujson?style=flat-square)](https://goreportcard.com/report/github.com/tomtwinkle/utfbomremover)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/tomtwinkle/utfbomremover/Build%20Check)

## Overview
`transform.Transformer` to remove Unicode BOM (Byte Order Mark).

UnicodeのBOM(Byte Order Mark)を削除する`transform.Transformer`。

## Usage

```golang
const base = []byte("一二三四五六七八九十拾壱")
msg := append([]byte{0xEF, 0xBB, 0xBF}, base...)

// true
fmt.Println(utfbomremover.ISUTF8BOM(msg))

var buf bytes.Buffer
w := transform.NewWriter(&buf, utfbomremover.NewTransformer())
if _, err := w.Write(msg); err != nil {
    panic(err)
}
if err := w.Close(); err != nil {
    panic(err)
}

// false
fmt.Println(utfbomremover.ISUTF8BOM(buf.Bytes()))
```