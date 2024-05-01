package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type customResponseWriter struct {
	http.ResponseWriter
	CompressWriter io.Writer
	buf            *bytes.Buffer
	statusCode     int
	responseSize   int
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *customResponseWriter) Write(data []byte) (int, error) {
	var size int
	var err error
	w.buf.Write(data)
	if w.CompressWriter != nil {
		size, err = w.CompressWriter.Write(data)
	} else {
		size, err = w.ResponseWriter.Write(data)
	}
	w.responseSize += size
	return size, err
}

func (w *customResponseWriter) ReadAll() ([]byte, error) {
	return w.buf.Bytes(), nil
}

func CustomizeResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crw := &customResponseWriter{
			ResponseWriter: w,
			buf:            bytes.NewBuffer(nil),
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
			defer gz.Close()
			crw.CompressWriter = gz
		}
		next.ServeHTTP(crw, r)
	})
}
