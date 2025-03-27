package lib

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

			//close writer before accessing buf.Bytes()
			if err := gzipWriter.Close(); err != nil {
				return "", err
			}

			fmt.Printf("Compressed size: %d bytes\n", buf.Len())
			// os.WriteFile("compressed_data.gz", buf.Bytes(), 0644)
			return bytesToBase64(buf.Bytes()), nil
		},
		// Decompress: func(s string) (T, error) {
		// 	var out T

		// 	b, err := base64ToBytes(s)
		// 	if err != nil {
		// 		return out, err
		// 	}

		// 	gzipReader, err := gzip.NewReader(bytes.NewReader(b))
		// 	if err != nil {
		// 		return out, err
		// 	}
		// 	defer gzipReader.Close()

		// 	jsnDecoder := json.NewDecoder(gzipReader)
		// 	if err := jsnDecoder.Decode(&out); err != nil {
		// 		return out, err
		// 	}
		// 	return out, nil
		// },
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

			raw, err := io.ReadAll(gzipReader)
			if err != nil {
				return out, fmt.Errorf("reading gzip stream failed: %w", err)
			}

			os.WriteFile("data.json", raw, 0644)
			fmt.Println("RAW JSON DUMP (first 1k chars):", string(raw[:1000])) // or dump to file

			err = json.Unmarshal(raw, &out)
			if err != nil {
				return out, fmt.Errorf("unmarshal error: %w", err)
			}

			// jsnDecoder := json.NewDecoder(gzipReader)
			// if err := jsnDecoder.Decode(&out); err != nil {
			// 	return out, err
			// }
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
