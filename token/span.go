package token

type Span struct {
	Offset int
	Length int
}

func (s Span) Slice(str string) string {
	return str[s.Offset:s.Length]
}
