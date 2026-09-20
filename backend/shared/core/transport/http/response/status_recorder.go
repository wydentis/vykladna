package core_http_response

import "net/http"

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode  int
	wroteHeader bool
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		wroteHeader:    false,
	}
}

func (w *StatusRecorder) WriteHeader(code int) {
	if code >= 100 && code < 200 && code != http.StatusSwitchingProtocols {
		w.ResponseWriter.WriteHeader(code)
		return
	}
	if w.wroteHeader {
		return
	}

	w.wroteHeader = true
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *StatusRecorder) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(b)
}

func (w *StatusRecorder) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
