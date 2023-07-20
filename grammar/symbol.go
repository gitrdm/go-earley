package grammar

type Symbol interface {
	symbol()
}

type SymbolImpl struct{}

func (s *SymbolImpl) symbol() {}
