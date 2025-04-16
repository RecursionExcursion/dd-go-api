package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

type SerializeFn[R any] = func([]byte) (R, error)
type DeserializeFn[R any] = func(R) ([]byte, error)

type Compressor[T any, R any] struct {
	compressor dataCompressor[T]
	Compress   func(T) (R, error)
	Decompress func(R) (T, error)
}

func NewCompressor[T any, R any](
	serializeToRFn SerializeFn[R],
	deserializeFromRFn DeserializeFn[R],
) Compressor[T, R] {
	c := Compressor[T, R]{
		compressor: NewGzipper[T](),
	}

	/* Define fns  */

	c.Compress = func(t T) (R, error) {
		b, err := c.compressor.Compress(t)
		if err != nil {
			var r R
			return r, err
		}
		return serializeToRFn(b)
	}

	c.Decompress = func(r R) (T, error) {
		b, err := deserializeFromRFn(r)

		if err != nil {
			var t T
			return t, err
		}

		return c.compressor.Decompress(b)
	}

	return c
}

type dataCompressor[T any] struct {
	Compress   func(data T) ([]byte, error)
	Decompress func(b []byte) (T, error)
}

func NewGzipper[T any]() dataCompressor[T] {
	return dataCompressor[T]{
		Compress: func(data T) ([]byte, error) {
			var buf bytes.Buffer
			var emptyBuffer = []byte{}

			gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
			if err != nil {
				return emptyBuffer, err
			}

			jsnEncoder := json.NewEncoder(gzipWriter)
			if err := jsnEncoder.Encode(data); err != nil {
				return emptyBuffer, err
			}

			if err := gzipWriter.Close(); err != nil {
				return emptyBuffer, err
			}

			return buf.Bytes(), nil
		},
		Decompress: func(b []byte) (T, error) {
			var out T

			gzipReader, err := gzip.NewReader(bytes.NewReader(b))
			if err != nil {
				return out, err
			}
			defer gzipReader.Close()

			jsnDecoder := json.NewDecoder(gzipReader)
			if err := jsnDecoder.Decode(&out); err != nil {
				return out, err
			}
			return out, nil
		},
	}
}
