package valueobjects

type Status string

const (
	Created   Status = "created"
	Canceled  Status = "canceled"
	Delivers  Status = "delivers"
	Completed Status = "completed"
)
