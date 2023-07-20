package token

// Accumulator is used to accumulate runes for span creation. It is meant to be the single location of input memory to avoid duplication.
type Accumulator interface {
	Accumulate(ch rune)
	Span(start, length int) AccumulatorSpan
	Length() int
	Capture() []rune
	RuneAt(index int) rune
}

// NewAccumulator creates a new accumulator for storing input text
func NewAccumulator() Accumulator {
	return &accumulator{
		builder: []rune{},
	}
}

type accumulator struct {
	builder []rune
}

func (a *accumulator) Accumulate(ch rune) {
	a.builder = append(a.builder, ch)
}

func (a *accumulator) Length() int {
	return len(a.builder)
}

func (a *accumulator) Capture() []rune {
	return a.builder
}

func (a *accumulator) String() string {
	return string(a.builder)
}

func (a *accumulator) RuneAt(index int) rune {
	return a.builder[index]
}

func (a *accumulator) Span(start, length int) AccumulatorSpan {
	return &accumulatorSpan{
		accumulator: a,
		start:       start,
		length:      length,
	}
}
