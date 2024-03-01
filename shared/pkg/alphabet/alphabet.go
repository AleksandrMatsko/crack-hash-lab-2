package alphabet

import "strings"

type Alphabet struct {
	runes []rune
}

func InitAlphabet(runes []rune) Alphabet {
	return Alphabet{runes: runes}
}

func (a *Alphabet) Length() int {
	return len(a.runes)
}

func (a *Alphabet) GetWord(ids []uint64) string {
	builder := strings.Builder{}
	for _, v := range ids {
		builder.WriteRune(a.runes[v])
	}
	return builder.String()
}

func (a *Alphabet) ToOneLine() string {
	return string(a.runes)
}
