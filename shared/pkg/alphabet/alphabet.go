package alphabet

import (
	"strings"
)

type Alphabet struct {
	Runes []rune
}

func InitAlphabet(runes []rune) Alphabet {
	return Alphabet{Runes: runes}
}

func (a *Alphabet) Length() int {
	return len(a.Runes)
}

func (a *Alphabet) GetWord(ids []uint64) string {
	builder := strings.Builder{}
	for _, v := range ids {
		builder.WriteRune(a.Runes[v])
	}
	return builder.String()
}

func (a *Alphabet) ToOneLine() string {
	return string(a.Runes)
}
