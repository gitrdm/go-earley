package token

// Lexeme is a mutable token
type Lexeme interface {
	Scan(ch rune) bool
	Accepted() bool
}
