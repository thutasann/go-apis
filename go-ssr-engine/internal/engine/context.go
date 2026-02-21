package engine

// RenderContext holds runtime values for variables.
//
// IMPORTANT:
// - No maps (avoid hashing cost in hot path)
// - Indexes slice lookup only
// - Caller must ensure len(Values) >= VarCount
type RenderContext struct {
	Values [][]byte
}

// NewRenderContext builds RenderContext from map[string]string
// using variable order defined in template (varName -> index map)
func NewRenderContext(varMap map[string]string, tpl *Template, varNameToIndex map[string]uint16) RenderContext {
	values := make([][]byte, tpl.VarCount)
	for name, idx := range varNameToIndex {
		if val, ok := varMap[name]; ok {
			values[idx] = []byte(val)
		} else {
			values[idx] = []byte{} // empty string if missing
		}
	}
	return RenderContext{Values: values}
}
