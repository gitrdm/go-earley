package token

// ReadonlyAccumulatorSpan is a span of text that can not be modified
type ReadonlyAccumulatorSpan interface {
	Start() int
	Length() int
	Capture() []rune
	RuneAt(index int) rune
	Peek() rune
}

// EOF is the end of the stream
const EOF = rune(-1)

// AccumulatorSpan represnets a read/write span and is used to grow the span
type AccumulatorSpan interface {
	ReadonlyAccumulatorSpan
	Grow() bool
}

type accumulatorSpan struct {
	accumulator Accumulator
	start       int
	length      int
}

func (s *accumulatorSpan) Start() int {
	return s.start
}

func (s *accumulatorSpan) Length() int {
	return s.length
}

func (s *accumulatorSpan) Capture() []rune {
	return s.accumulator.Capture()[s.start:s.length]
}

func (s *accumulatorSpan) Grow() bool {
	if s.Peek() == EOF {
		return false
	}
	s.length++
	return true
}

func (s *accumulatorSpan) Peek() rune {
	if s.length >= s.accumulator.Length() {
		return EOF
	}
	return s.accumulator.RuneAt(s.length)
}

func (s *accumulatorSpan) RuneAt(index int) rune {
	return s.accumulator.RuneAt(index)
}
