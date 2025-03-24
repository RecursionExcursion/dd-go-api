package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// TODO consider make this take in reader and writer interfaces instead of strings
type genericGzipCompressor[T any] struct {
	Compress   func(data T) (string, error)
	Decompress func(s string) (T, error)
}

func GzipCompressor[T any]() genericGzipCompressor[T] {
	return genericGzipCompressor[T]{
		Compress: func(data T) (string, error) {
			var buf bytes.Buffer
			gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
			if err != nil {
				return "", err
			}

			//stream json directly into gzipWriter
			jsnEncoder := json.NewEncoder(gzipWriter)
			if err := jsnEncoder.Encode(data); err != nil {
				return "", err
			}

			fmt.Printf("Compressed size: %d bytes\n", buf.Len())

			//close writer befor accessing buf.Bytes()
			if err := gzipWriter.Close(); err != nil {
				return "", err
			}

			return bytesToBase64(buf.Bytes()), nil
		},
		Decompress: func(s string) (T, error) {
			var out T

			b, err := base64ToBytes(s)
			if err != nil {
				return out, err
			}

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

func bytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func base64ToBytes(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
