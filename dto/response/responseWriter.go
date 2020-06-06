package response

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SnmpResponseWriter struct {
	writer http.ResponseWriter
	logger *logrus.Logger
}

func (r *SnmpResponseWriter) WithResponse(res interface{}) {
	if err := json.NewEncoder(r.writer).Encode(&res); err != nil {
		r.logger.Println("Some error while encoding response: ", err.Error())
	}
}

func (r *SnmpResponseWriter) WithStatus(status int) *SnmpResponseWriter {
	r.writer.WriteHeader(status)
	return r
}

func New(rw http.ResponseWriter, logger *logrus.Logger) *SnmpResponseWriter {
	w := SnmpResponseWriter{writer: rw, logger: logger}
	return &w
}
