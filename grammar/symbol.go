package grammar

type Symbol interface {
	symbol()
	String() string
}

type SymbolImpl struct{}

func (s *SymbolImpl) symbol() {}
