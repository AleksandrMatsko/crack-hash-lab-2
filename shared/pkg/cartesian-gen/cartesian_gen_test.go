package cartesian_gen

import "testing"

func TestCartesianGenerator_Product(t *testing.T) {
	dims := []uint64{1, 2, 3}
	expected := [][]uint64{
		{0, 0, 0},
		{0, 0, 1},
		{0, 0, 2},
		{0, 1, 0},
		{0, 1, 1},
		{0, 1, 2},
	}
	t.Logf("len(expected) = %v, len(expected[0] = %v)", len(expected), len(expected[0]))

	gen := NewCartesianGenerator(dims)
	for i := 0; i < len(expected); i++ {
		pr := gen.Product()
		t.Logf("expected[%v] = %v got = %v", i, expected[i], pr)
		for j := 0; j < len(expected[i]); j++ {
			if expected[i][j] != pr[j] {
				println(expected[i])
				println(pr)
				t.Fatalf("bad value on step %v index %v", i, j)
			}
		}
	}

	t.Logf("testing generator with same limit")
	gen = NewCartesianGenerator(dims).Limit(uint64(len(expected)))
	for i := 0; i < len(expected); i++ {
		pr := gen.Product()
		t.Logf("expected[%v] = %v got = %v", i, expected[i], pr)
		for j := 0; j < len(expected[i]); j++ {
			if expected[i][j] != pr[j] {
				println(expected[i])
				println(pr)
				t.Fatalf("bad value on step %v index %v", i, j)
			}
		}
	}
}

func TestCartesianGenerator_withOneDim(t *testing.T) {
	dims := []uint64{36}
	expected := make([][]uint64, 0)
	for i := 0; i < int(dims[0]); i++ {
		expected = append(expected, []uint64{uint64(i)})
	}

	gen := NewCartesianGenerator(dims)
	for i := 0; i < len(expected); i++ {
		got := gen.Product()
		t.Logf("expected[%v] = %v got = %v", i, expected[i], got)
		if expected[i][0] != got[0] {
			t.Fatalf("bad values")
		}
	}
}

func TestCartesianGenerator_Skip(t *testing.T) {
	dims := []uint64{1, 2, 3}
	expected := [][]uint64{
		{0, 0, 0},
		{0, 0, 1},
		{0, 0, 2},
		{0, 1, 0},
		{0, 1, 1},
		{0, 1, 2},
	}
	t.Logf("len(expected) = %v, len(expected[0] = %v)", len(expected), len(expected[0]))
	var skip uint64 = 1

	gen := NewCartesianGenerator(dims).Skip(skip)
	for i := int(skip); i < len(expected); i++ {
		pr := gen.Product()
		t.Logf("expected[%v] = %v got = %v", i, expected[i], pr)
		for j := 0; j < len(expected[i]); j++ {
			if expected[i][j] != pr[j] {
				println(expected[i])
				println(pr)
				t.Fatalf("bad value on step %v index %v", i, j)
			}
		}
	}
}

func TestCartesianGenerator_Limit(t *testing.T) {
	dims := []uint64{1, 2, 3}
	limit := 3

	gen := NewCartesianGenerator(dims).Limit(uint64(limit))
	i := 0
	for gen.HasNext() {
		gen.Product()
		i += 1
	}
	if i != limit {
		t.Fatalf("expected %v iterations, got %v iterations", limit, i)
	}
}
