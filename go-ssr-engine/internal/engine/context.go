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
