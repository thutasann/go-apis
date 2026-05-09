package job

// Priority represents job importance level.
// We keep it small (2 levels) for predictable behavior
type Priority int

const (
	// High Priority jobs are processed first.
	High Priority = iota

	// Normal Priority jobs are processed after High
	Normal
)

// String makes logs readable.
func (p Priority) String() string {
	switch p {
	case High:
		return "HIGH"
	case Normal:
		return "NORMAL"
	default:
		return "UNKNOWN"
	}
}
