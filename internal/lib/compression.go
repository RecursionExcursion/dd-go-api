package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

//TODO Write doc

type Algorithms = struct {
	Writer func(w io.Writer) (io.WriteCloser, error)
	Reader func(r io.Reader) (io.Reader, error)
}

type Codec[R any] struct {
	Encode func([]byte) (R, error)
	Decode func(R) ([]byte, error)
}

type Compressor[T any, R any] struct {
	Algorithms Algorithms
	Compress   func(T) (R, error)
	Decompress func(R) (T, error)
}

func GzipCompressor[T any, R any](codec Codec[R]) Compressor[T, R] {

	c := Compressor[T, R]{
		Algorithms: Algorithms{
			Writer: func(w io.Writer) (io.WriteCloser, error) {
				return gzip.NewWriterLevel(w, gzip.BestCompression)
			},
			Reader: func(r io.Reader) (io.Reader, error) {
				return gzip.NewReader(r)
			},
		},
	}
	c.assignCodec(codec)
	return c
}

func CustomCompressor[T any, R any](algos Algorithms, codec Codec[R]) Compressor[T, R] {
	c := Compressor[T, R]{
		Algorithms: algos,
	}
	c.assignCodec(codec)
	return c
}

func (c *Compressor[T, R]) assignCodec(codec Codec[R]) {
	c.Compress = func(t T) (R, error) {
		b, err := c.compressT(t)
		if err != nil {
			var r R
			return r, err
		}
		return codec.Encode(b)
	}

	c.Decompress = func(r R) (T, error) {
		b, err := codec.Decode(r)

		if err != nil {
			var t T
			return t, err
		}

		return c.decompressT(b)
	}
}

func (c *Compressor[T, R]) compressT(data T) ([]byte, error) {
	var buf bytes.Buffer
	var emptyBuffer = []byte{}

	w, err := c.Algorithms.Writer(&buf)
	if err != nil {
		return emptyBuffer, err
	}

	jsnEncoder := json.NewEncoder(w)
	if err := jsnEncoder.Encode(data); err != nil {
		return emptyBuffer, err
	}

	if err := w.Close(); err != nil {
		return emptyBuffer, err
	}

	return buf.Bytes(), nil
}

func (c *Compressor[T, R]) decompressT(b []byte) (T, error) {
	var out T

	r, err := c.Algorithms.Reader(bytes.NewReader(b))
	if err != nil {
		return out, err
	}
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	jsnDecoder := json.NewDecoder(r)
	if err := jsnDecoder.Decode(&out); err != nil {
		return out, err
	}
	return out, nil
}
