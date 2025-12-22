package sse

import (
	"encoding/json"
	"net/http"
)

var headers = map[string]string{
	"Access-Control-Allow-Origin":   "*",
	"Access-Control-Expose-Headers": "Content-Type",
	"Content-Type":                  "text/event-stream",
	"Cache-Control":                 "no-cache",
	"Connection":                    "keep-alive",
}

type SSEWriter struct {
	w   http.ResponseWriter
	enc *json.Encoder
}

func NewSSEWriter(w http.ResponseWriter) *SSEWriter {
	return &SSEWriter{w: w, enc: json.NewEncoder(w)}
}

func (s *SSEWriter) WriteHeaders() {
	for key, value := range headers {
		s.w.Header().Set(key, value)
	}
}

func (s *SSEWriter) WriteEvent(evtname string, payload any) {
	s.w.Write([]byte("event:" + evtname + "\ndata:"))
	s.enc.Encode(payload)
	s.w.Write([]byte("\n\n"))

	if f, ok := s.w.(http.Flusher); ok {
		f.Flush()
	}
}
