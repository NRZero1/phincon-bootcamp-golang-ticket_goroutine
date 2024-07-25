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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWrap := responseWriterWrap{
			ResponseWriter: w,
			statusCode: http.StatusOK,
			statusDesc: http.StatusText(http.StatusOK),
		}
		start := time.Now()
		next.ServeHTTP(responseWrap, r)
		log.Info().Msg(fmt.Sprintf("%d %s %s %s %d", responseWrap.statusCode, responseWrap.statusDesc, r.Method, r.URL.Path, time.Since(start)))
		// fmt.Printf("%s %s %d", r.Method, r.URL.Path, time.Since(start))
	})
}
