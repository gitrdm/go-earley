package capture

type Builder struct {
	slice []rune
}

func FromSlice(slice ...rune) *Builder {
	return &Builder{
		slice: slice,
	}
}

// RuneAt implements Capture
func (b *Builder) RuneAt(index int) rune {
	return b.slice[index]
}

// Len implements Capture
func (b *Builder) Len() int {
	return len(b.slice)
}

// Offset implements Capture
func (*Builder) Offset() int {
	return 0
}

// Append implements Appendable
func (b *Builder) Append(r rune) {
	b.slice = append(b.slice, r)
}
