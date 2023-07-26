package capture

type Capture interface {
	Offset() int
	Len() int
	RuneAt(index int) rune
}

type Appendable interface {
	Append(r rune)
}

type MutableOffset interface {
	SetOffset(i int)
}

type MutableLen interface {
	SetLen(i int)
}
