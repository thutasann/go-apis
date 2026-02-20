package engine

// Compile parses template bytes and converts them
// into a slice of Instructions.
//
// Supported syntax:
//
//	{{variable}}
//
// Design goals:
// - Single pass scan
// - No regex
// - No reflection
// - Minimal allocations
func Compile(input []byte) (*Template, error) {
	var instructions []Instruction
	var vars [][]byte

	i := 0
	for i < len(input) {

		// Detect variable start
		if i+1 < len(input) && input[i] == '{' && input[i+1] == '{' {

			// Flush any previous static text
			start := i
			if start > 0 {
				text := input[:start]
				if len(text) > 0 {
					instructions = append(instructions, Instruction{
						Op:   OpText,
						Data: append([]byte(nil), text...), // copy to isolate
					})
				}
			}

			// Move past {{
			i += 2
			varStart := i

			// Find closing }}
			for i+1 < len(input) && !(input[i] == '}' && input[i+1] == '}') {
				i++
			}

			varName := input[varStart:i]

			// Assign index based on order of appearance
			idx := uint16(len(vars))
			vars = append(vars, append([]byte(nil), varName...))

			instructions = append(instructions, Instruction{
				Op:  OpVar,
				Idx: idx,
			})

			// Move past }}
			i += 2

			// Truncate input to remaining slice
			input = input[i:]
			i = 0
			continue
		}

		i++
	}

	// Remaining static text
	if len(input) > 0 {
		instructions = append(instructions, Instruction{
			Op:   OpText,
			Data: append([]byte(nil), input...),
		})
	}

	return &Template{
		Instructions: instructions,
		VarCount:     uint16(len(vars)),
	}, nil
}
