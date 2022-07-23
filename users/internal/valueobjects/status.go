package valueobjects

type Status string

const (
	Active  Status = "active"
	Frozen  Status = "frozen"
	Blocked Status = "blocked"
)
