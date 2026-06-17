package models

// DriverStanding represents a driver's season total.
type DriverStanding struct {
	Position   int
	DriverName string
	Team       string
	Points     float64
}

// TeamStanding represents a constructor/team's season total.
type TeamStanding struct {
	Position int
	TeamName string
	Points   float64
}

// PastLeader represents a leader entry during a grand prix (simple form).
// For now it's just the driver's name and laps led could be added later.
type PastLeader struct {
	Name string
}
