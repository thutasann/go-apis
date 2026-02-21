package engine

import "bytes"

// Compile parses template and builds instruction list.
//
// Supported syntax:
//
//	{{var}}
//
// No regex.
// Single pass.
// Minimal allocations.
func Compile(input []byte) (*Template, error) {
	var instructions []Instruction
	varNames := make(map[string]uint16)

	i := 0
	start := 0

	for i < len(input) {

		// detect {{
		if i+1 < len(input) && input[i] == '{' && input[i+1] == '{' {

			// flush text before {{
			if start < i {
				instructions = append(instructions, Instruction{
					Op:   OpText,
					Data: append([]byte(nil), input[start:i]...),
				})
			}

			i += 2
			varStart := i

			// find }}
			for i+1 < len(input) && !(input[i] == '}' && input[i+1] == '}') {
				i++
			}

			varName := string(bytes.TrimSpace(input[varStart:i]))

			// assign stable index per variable name
			idx, exists := varNames[varName]
			if !exists {
				idx = uint16(len(varNames))
				varNames[varName] = idx
			}

			instructions = append(instructions, Instruction{
				Op:  OpVar,
				Idx: idx,
			})

			i += 2
			start = i
			continue
		}

		i++
	}

	// remaining text
	if start < len(input) {
		instructions = append(instructions, Instruction{
			Op:   OpText,
			Data: append([]byte(nil), input[start:]...),
		})
	}

	return &Template{
		Instructions: instructions,
		VarCount:     uint16(len(varNames)),
	}, nil
}
