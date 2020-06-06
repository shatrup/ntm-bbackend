package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type ApiAuthMiddleware struct {
	Logger *logrus.Logger
}

func (am ApiAuthMiddleware) ApiAuth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessTokenFromHeader := r.Header.Get("NTM-AUTH-KEY")
	if accessTokenFromHeader == os.Getenv("NTM_API_KEY") {
		next(rw, r)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}
