package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type responseWriterWrap struct {
	http.ResponseWriter
	statusCode int
	statusDesc string
}

func (rw *responseWriterWrap) WriteHeader(code int) {
	rw.statusCode = code
	rw.statusDesc = http.StatusText(code)
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWrap := responseWriterWrap{
			ResponseWriter: w,
			statusCode: http.StatusOK,
			statusDesc: http.StatusText(http.StatusOK),
		}
		start := time.Now()
		next.ServeHTTP(&responseWrap, r)
		log.Info().Msg(fmt.Sprintf("%d %s process time: %d ns %s %s", responseWrap.statusCode, responseWrap.statusDesc, time.Since(start), r.Method, r.URL.Path))
		// log.Info().Msg(fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, time.Since(start)))
		// fmt.Printf("%s %s %d", r.Method, r.URL.Path, time.Since(start))
	})
}
