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
	underlying *http.ResponseWriter
	enc        *json.Encoder
}

func NewSSEWriter(underlying *http.ResponseWriter) *SSEWriter {
	enc := json.NewEncoder(*underlying)

	return &SSEWriter{
		underlying,
		enc,
	}
}

func (s *SSEWriter) WriteHeaders() {
	for key, value := range headers {
		(*s.underlying).Header().Set(key, value)
	}
}

func (s *SSEWriter) WriteEvent(evtname string, payload any) {
	(*s.underlying).Write([]byte("event:" + evtname + "\ndata:"))
	s.enc.Encode(payload)
	(*s.underlying).Write([]byte("\n\n"))
	(*s.underlying).(http.Flusher).Flush()
}
