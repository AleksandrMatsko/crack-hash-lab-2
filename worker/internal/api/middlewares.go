package api

import (
	"context"
	"distributed.systems.labs/worker/internal/notify"
	"net/http"
)

type notifierKeyType int

var notifierKey notifierKeyType

func getNotifier(ctx context.Context) notify.Notifier {
	return ctx.Value(notifierKey).(notify.Notifier)
}

type notifierMiddleware struct {
	mn notify.Notifier
}

func (m *notifierMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), notifierKey, m.mn)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
