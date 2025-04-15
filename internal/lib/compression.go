package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

type DataCompressor[T any] struct {
	Compress   func(data T) ([]byte, error)
	Decompress func(b []byte) (T, error)
}

func NewGzipper[T any]() DataCompressor[T] {
	return DataCompressor[T]{
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
