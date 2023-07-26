// The capture package is responsible for rune slice pointer management
// When a parser is scanning input, additional characters will be added to the
// rune slice. Each token has a pointer to a slice of that rune slice but appending
// would cause it to lose track of the original pointer.

// Captures provide pointer wrapping of the original base rune which forms a tree.
// When the original underlying pointer is appended to, only the capture at the base of the tree is updated
// All other slices into the original capture point to the base capture and not the underlying array

// The main purpose of this package is to save each token retaining a copy of a string. This parser can parse ambiguity
// so duplicates of the same underlying string often exist

// A secondary purpose of this package is to remove the need to read the entire file into memory at once.
package capture
