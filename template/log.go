package {{Name}}

import (
	"io/ioutil"
	stdlog "log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-pg/pg"
)

var log = stdlog.New(os.Stderr, "", stdlog.Ldate|stdlog.Ltime|stdlog.Lshortfile)

func DisableLog() {
	log.SetOutput(ioutil.Discard)
}

type Logger struct{}

func (l Logger) BeforeQuery(q *pg.QueryEvent) {
	q.Data["st"] = time.Now()
}

func (l Logger) AfterQuery(q *pg.QueryEvent) {
	query, err := q.FormattedQuery()
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("%s took %s", query, time.Since(q.Data["st"].(time.Time)))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (l Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				stack := make([]byte, 1024*8)
				log.Print(string(stack[:runtime.Stack(stack, false)]))
			}
		}()
		st := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(w, req)
		log.Printf("%v %s %s %s", lrw.statusCode, req.Method, req.URL.Path, time.Since(st))
	})
}
