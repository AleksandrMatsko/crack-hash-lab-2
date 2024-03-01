package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type loggerKeyType int

var loggerKey loggerKeyType

func GetLogger(ctx context.Context) *log.Logger {
	return ctx.Value(loggerKey).(*log.Logger)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defaultLog := log.Default()
		customLogger := log.New(
			defaultLog.Writer(),
			fmt.Sprintf("%s %s %s ", r.Method, r.URL, r.Host),
			defaultLog.Flags()|log.Lmsgprefix)
		ctx := context.WithValue(r.Context(), loggerKey, customLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
