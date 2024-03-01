package cartesian_gen

type CartesianGenerator struct {
	currentWord []uint64
	dims        []uint64
	hasLimit    bool
	limit       uint64
}

func NewCartesianGenerator(dims []uint64) *CartesianGenerator {
	return &CartesianGenerator{
		currentWord: make([]uint64, len(dims)),
		dims:        dims,
		hasLimit:    false,
		limit:       1,
	}
}

func (g *CartesianGenerator) Product() []uint64 {
	ret := make([]uint64, len(g.currentWord))
	copy(ret, g.currentWord)

	var carry uint64 = 1

	for i := len(g.currentWord) - 1; i >= 0 && carry > 0 && g.limit > 0; i-- {
		g.currentWord[i] += carry
		carry = g.currentWord[i] / g.dims[i]
		g.currentWord[i] %= g.dims[i]
	}
	if g.hasLimit {
		g.limit -= 1
	}
	if carry == 1 || g.limit == 0 {
		g.currentWord = make([]uint64, 0)
	}
	return ret
}

func (g *CartesianGenerator) Skip(val uint64) *CartesianGenerator {
	if val == 0 {
		return g
	}
	word := make([]uint64, len(g.currentWord))
	if len(word) == 0 {
		return g
	}
	carry := val
	for i := len(word) - 1; i >= 0; i-- {
		word[i] += carry
		carry = word[i] / g.dims[i]
		word[i] %= g.dims[i]
	}
	if carry > 0 {
		g.currentWord = make([]uint64, 0)
		return g
	}
	carry = 0
	for i := len(g.currentWord) - 1; i >= 0; i-- {
		g.currentWord[i] += word[i] + carry
		carry = g.currentWord[i] / g.dims[i]
		g.currentWord[i] %= g.dims[i]
	}
	if carry > 0 {
		g.currentWord = make([]uint64, 0)
	}
	return g
}

func (g *CartesianGenerator) Limit(limit uint64) *CartesianGenerator {
	g.hasLimit = true
	g.limit = limit + 1
	return g
}

func (g *CartesianGenerator) HasNext() bool {
	return len(g.currentWord) != 0
}
