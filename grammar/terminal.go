package grammar

type Terminal interface {
	Symbol
	IsMatch(ch rune) bool
}
