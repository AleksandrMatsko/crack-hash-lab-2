package middlewares

import (
	"context"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/shared/pkg/alphabet"
	"net/http"
)

type storageKeyType int

var storageKey storageKeyType

func GetStorage(ctx context.Context) storage.Storage {
	return ctx.Value(storageKey).(storage.Storage)
}

type StorageMiddleware struct {
	S storage.Storage
}

func (s *StorageMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), storageKey, s.S)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type alphabetKeyType int

var alphabetKey alphabetKeyType

func GetAlphabet(ctx context.Context) alphabet.Alphabet {
	return ctx.Value(alphabetKey).(alphabet.Alphabet)
}

type AlphabetMiddleware struct {
	A alphabet.Alphabet
}

func (a *AlphabetMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), alphabetKey, a.A)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
