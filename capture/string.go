package capture

func FromString(str string) Capture {
	return FromSlice([]rune(str)...)
}
