package tasks

import (
	"math/rand"
	"testing"
)

func TestPrivate_pow(t *testing.T) {
	base := uint64(rand.Intn(10))
	exp := uint64(rand.Intn(10))

	var expected uint64 = 1
	for i := 0; uint64(i) < exp; i++ {
		expected *= base
	}

	powRes := pow[uint64](base, exp)
	if powRes != expected {
		t.Fatalf("base = %v exp = %v expected %v got %v", base, exp, expected, powRes)
	}
}

func TestPrivate_calcTotalWordsCount(t *testing.T) {
	lenAlphabet := uint64(rand.Intn(100))
	maxLength := uint64(rand.Intn(4))

	var expected uint64 = 0
	deg := lenAlphabet
	for i := 0; uint64(i) < maxLength; i++ {
		expected += deg
		deg *= lenAlphabet
	}

	got := calcTotalWordsCount(lenAlphabet, maxLength)
	if expected != got {
		t.Fatalf("lenAlphabet = %v maxLength = %v expected %v got %v", lenAlphabet, maxLength, expected, got)
	}
}
