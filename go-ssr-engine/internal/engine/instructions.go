package engine

// OpCode represents low-level instruction operation.
// Using uint8 keeps struct small and cache friendly
type OpCode uint8

const (
	// OpText writes static text directly to output
	OpText OpCode = iota

	// OpVar writes a variable value by index lookup.
	// Index referes to position in RenderContext.Values.
	OpVar
)

// Instruction is a single compiled template operation.
//
// Design goals:
// - No interfaces (avoid dynamic dispatch)
// - No reflection
// - Predictable memory layout
// - Branch-friendly in hot loop
//
// Memory layout:
//
//	Op   -> 1 byte
//	Data -> slice header (24 bytes)
//	Idx  -> 2 bytes
//
// Total small enough to stay cache friendly.
type Instruction struct {
	Op   OpCode // operation type
	Data []byte // static text (used when OpText)
	Idx  uint16 // variable index (used when OpVar)
}
