package capture

type slice struct {
	parent Capture
	offset int
	len    int
}

// Len implements Capture.
func (s *slice) Len() int {
	return s.len
}

// SetLen implements MutableLen
func (s *slice) SetLen(length int) {
	s.len = length
}

// Offset implements Capture.
func (s *slice) Offset() int {
	return s.offset
}

// SetOffset implements MutableOffset
func (s *slice) SetOffset(offset int) {
	s.offset = offset
}

// RuneAt implements Capture.
func (s *slice) RuneAt(index int) rune {
	return s.parent.RuneAt(s.offset + index)
}

func Span(c Capture, index int) Capture {
	return &slice{
		parent: c,
		offset: c.Offset() + index,
		len:    c.Len() - index,
	}
}

func SpanLength(c Capture, start int, length int) Capture {
	return &slice{
		parent: c,
		offset: c.Offset() + start,
		len:    length,
	}
}
