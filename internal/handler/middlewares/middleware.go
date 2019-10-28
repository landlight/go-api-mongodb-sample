package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger logger
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path

		newRw := newResponseWriter(rw)
		next.ServeHTTP(newRw, r)

		clientIP := r.Header.Get("X-Forwarded-For")
		if userIP, ok := IPFromRequest(r); ok == nil {
			lastIP := userIP.String()
			if clientIP == "" {
				clientIP = lastIP
			} else {
				clientIP += fmt.Sprintf(", %s", lastIP)
			}
		}
		userAgent := r.Header.Get("User-Agent")
		contentType := r.Header.Get("Content-Type")
		hostname, _ := os.Hostname()
		method := r.Method
		statusCode := newRw.Status()
		duration := time.Since(start)
		normalLog := log.Fields{
			"statusCode":  statusCode,
			"latency":     duration.String(),
			"clientIP":    clientIP,
			"method":      method,
			"path":        path,
			"userAgent":   userAgent,
			"contentType": contentType,
			"hostname":    hostname,
		}
		//if path != "/api/healthz" {
		log.WithFields(normalLog).Info(fmt.Sprintf("latency: %s method: %s path: %s userAgent: %s", duration.String(), method, path, userAgent))
		// }
	})
}

type responseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseWriterModel struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(rw http.ResponseWriter) *responseWriterModel {
	nrw := &responseWriterModel{
		ResponseWriter: rw,
	}
	return nrw
}

func (rw *responseWriterModel) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriterModel) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	return size, err
}

func (rw *responseWriterModel) Status() int {
	return rw.status
}

// IPFromRequest get ip address
func IPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}
