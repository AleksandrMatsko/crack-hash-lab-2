package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetWorkers_WithBadFormat_Fails(t *testing.T) {
	ConfigureApp()

	err := os.Setenv(appEnvWorkersPrefix+"_LIST", "worker.")
	if err != nil {
		t.Fatalf("unexpected error while setting env: %s", err)
	}

	_, err = GetWorkers()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrBadWorkersListFormat) {
		t.Fatalf("bad error: got '%s', expected '%s'", err, ErrBadWorkersListFormat)
	}
}

func TestGetWorkers_WithExpectedFormat_Works(t *testing.T) {
	ConfigureApp()

	val := "worker-1:80:worker-2:90:worker-3:100"
	err := os.Setenv(appEnvWorkersPrefix+"_LIST", val)
	if err != nil {
		t.Fatalf("unexpected error while setting env: %s", err)
	}
	splited := strings.Split(val, ":")
	workers := make([]string, 0)
	for i := 0; i < len(splited); i += 2 {
		workers = append(workers, fmt.Sprintf("%s:%s", splited[i], splited[i+1]))
	}

	w, err := GetWorkers()
	if err != nil {
		t.Fatalf("unexpected error while getting workers: %s", err)
	}
	if len(w) != len(workers) {
		t.Fatalf("bad length of workers: got %d, expected %d", len(w), len(workers))
	}
	for i := range workers {
		if w[i] != workers[i] {
			t.Fatalf("worker[%d]: got '%s', expected '%s'", i, w[i], workers[i])
		}
	}
}
