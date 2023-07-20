package token

type Token interface {
	Reset()
	Scan() bool
	Accepted() bool
}
