package engine

// Template represents an immutable compiled template.
//
// Safe for concurrent use across goroutines
// Instructions slice must never be mutated after compile
type Template struct {
	Instructions []Instruction
	VarCount     uint16
}
