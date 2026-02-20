package engine

import "io"

// RenderTo renders template directly to io.Writer.
//
// Zero string allocations.
// No reflection.
// No map lookups.
// Only slice index operations.
//
// Hot path must remain extremely simple for branch prediction.
func (t *Template) RenderTo(w io.Writer, ctx *RenderContext) error {
	for _, inst := range t.Instructions {
		switch inst.Op {
		case OpText:
			// Static write
			if _, err := w.Write(inst.Data); err != nil {
				return err
			}

		case OpVar:
			// Indexed variable lookup
			if int(inst.Idx) >= len(ctx.Values) {
				continue
			}
			if _, err := w.Write(ctx.Values[inst.Idx]); err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderBytes renders template into pooled buffer and returns bytes.
//
// Caller must not retain returned slice after next pool reuse.
func (t *Template) RenderBytes(ctx *RenderContext) ([]byte, error) {
	buf := getBuffer()
	defer putBuffer(buf)

	if err := t.RenderTo(buf, ctx); err != nil {
		return nil, err
	}

	// copy before returning because buffer will be reused.
	out := append([]byte(nil), buf.Bytes()...)
	return out, nil
}
